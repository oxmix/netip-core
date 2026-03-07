package collector

import (
	"os"
	"regexp"
	"strconv"
	"time"
)

func (c *Collector) collectMem() {
	for range time.Tick(time.Second) {
		c.procMemInfoHandler("/proc/meminfo")
	}
}

func (c *Collector) procMemInfoHandler(file string) {
	memInfo, err := os.ReadFile(file)
	if err != nil {
		return
	}

	defer c.mu.Unlock()
	c.mu.Lock()

	re := regexp.MustCompile(`(.*?):.*?([0-9]+)`)
	for _, m := range re.FindAllStringSubmatch(string(memInfo), -1) {
		if len(m) < 2 {
			continue
		}
		val, err := strconv.Atoi(m[2])
		if err != nil {
			val = 0
		}

		switch m[1] {
		case "MemTotal":
			c.data.MemStats.MemTotal = val
		case "MemFree":
			c.data.MemStats.MemFree = val
		case "Buffers":
			c.data.MemStats.Buffers = val
		case "Cached":
			c.data.MemStats.Cached = val
		case "Slab":
			c.data.MemStats.Slab = val
		case "SwapTotal":
			c.data.MemStats.SwapTotal = val
		case "SwapFree":
			c.data.MemStats.SwapFree = val
		}
	}
}
