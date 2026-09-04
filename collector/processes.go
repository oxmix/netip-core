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

	reader := io.LimitReader(f, 4096)
	s, _ := io.ReadAll(reader)
	line := string(s)

	open := strings.IndexByte(line, '(')
	closeIdx := strings.LastIndexByte(line, ')')
	if open < 0 || closeIdx < 0 || closeIdx < open {
		return
	}

	p.Name = line[open+1 : closeIdx]

	rest := strings.TrimSpace(line[closeIdx+1:])
	fields := strings.Fields(rest)
	if len(fields) < 18 {
		return
	}
	p.State = fields[0]
	p.PPID, _ = strconv.Atoi(fields[1])
	p.Threads, _ = strconv.Atoi(fields[17])
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
