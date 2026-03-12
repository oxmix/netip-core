package info

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (i *Info) fillMemInfo() {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)
	var memKB, swapKB int64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			memKB, _ = strconv.ParseInt(fields[1], 10, 64)
		} else if strings.HasPrefix(line, "SwapTotal:") {
			fields := strings.Fields(line)
			swapKB, _ = strconv.ParseInt(fields[1], 10, 64)
		}
		if memKB > 0 && swapKB > 0 {
			break
		}
	}

	i.Data.Mem.RAM = float64(memKB) / 1024 / 1024
	i.Data.Mem.Swap = float64(swapKB) / 1024 / 1024
}
