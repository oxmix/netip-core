package collector

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var cpuPrev = map[int]*cpuCore{}

type cpuCore struct {
	idle  int
	total int
}

func (c *Collector) collectLoadAvg() {
	for {
		file, err := os.ReadFile("/proc/loadavg")
		if err != nil {
			return
		}
		c.data.LoadAvg = strings.SplitN(string(file), " ", 4)[0:3]
		time.Sleep(time.Second)
	}
}

func (c *Collector) collectCPU() {
	for {
		c.procStatHandler("/proc/stat")
		time.Sleep(time.Second)
	}
}

func (c *Collector) procStatHandler(file string) {
	stat, err := os.ReadFile(file)
	if err != nil {
		return
	}

	var cores []int
	total := 0
	re := regexp.MustCompile(`(cpu\d+).+`)
	for i, match := range re.FindAllString(string(stat), -1) {
		num := i + 1

		it := cpuIdleTotal(match)

		if _, ok := cpuPrev[num]; !ok {
			cpuPrev[num] = &cpuCore{0, 0}
		}

		diffIdle := it.idle - cpuPrev[num].idle
		diffTotal := it.total - cpuPrev[num].total

		if diffTotal <= 0 {
			diffIdle = it.idle
			diffTotal = it.total
		}

		cpuPrev[num].idle = it.idle
		cpuPrev[num].total = it.total

		core := (1000*(diffTotal-diffIdle)/diffTotal + 5) / 10
		total = total + core

		cores = append(cores, core)
	}

	defer c.mu.Unlock()
	c.mu.Lock()

	c.data.Time = time.Now().UTC()
	c.data.CPUStats.Cores = cores
	numCores := len(cores)
	if numCores > 0 {
		c.data.CPUStats.Avg = float32(total) / float32(numCores)
	}
}

func cpuIdleTotal(match string) cpuCore {
	line := strings.Split(match, " ")
	idle, _ := strconv.Atoi(line[4])
	total := 0
	for _, l := range line {
		i, err := strconv.Atoi(l)
		if err == nil {
			total += i
		}
	}
	return cpuCore{idle, total}
}
