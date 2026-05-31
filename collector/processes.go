package collector

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const pathProc = "/proc"

type Processes []Proc

func (p *Proc) path(stat string) string {
	return pathProc + "/" + strconv.Itoa(p.PID) + "/" + stat
}

func (p *Proc) quantityFd() {
	d, err := os.Open(p.path("fd"))
	if err != nil {
		return
	}
	defer func(d *os.File) {
		_ = d.Close()
	}(d)

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}

	p.FDs = len(names)
}

func (p *Proc) fillState() {
	f, err := os.Open(p.path("stat"))
	if err != nil {
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := io.LimitReader(f, 1024)
	s, _ := io.ReadAll(reader)

	sn := strings.SplitN(string(s), " ", 21)
	if len(sn) < 19 {
		return
	}
	p.Name = sn[1][1 : len(sn[1])-1]
	p.State = sn[2]
	p.PPID, _ = strconv.Atoi(sn[3])
	p.Threads, _ = strconv.Atoi(sn[19])
}

func (c *Collector) collectProc() {
	for range time.Tick(30 * time.Second) {
		d, err := os.Open(pathProc)
		if err != nil {
			return
		}
		names, err := d.Readdirnames(-1)
		if err != nil {
			return
		}
		_ = d.Close()

		p := Processes{}
		for _, n := range names {
			pid, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				continue
			}
			pc := Proc{PID: int(pid)}
			pc.quantityFd()
			pc.fillState()
			p = append(p, pc)
		}
		c.ChanProcesses <- &p
	}
}
