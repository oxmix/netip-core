package info

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	reTwoColumns = regexp.MustCompile("\t+: ")
	reExtraSpace = regexp.MustCompile(" +")
	reCacheSize  = regexp.MustCompile(`^(\d+) KB$`)
)

func (i *Info) fillCPUInfo() {
	i.Data.CPU.Threads = uint(runtime.NumCPU())

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	cpu := make(map[string]bool)
	core := make(map[string]bool)

	var cpuID string

	s := bufio.NewScanner(f)
	for s.Scan() {
		if sl := reTwoColumns.Split(s.Text(), 2); sl != nil {
			switch sl[0] {
			case "physical id":
				cpuID = sl[1]
				cpu[cpuID] = true
			case "core id":
				coreID := fmt.Sprintf("%s/%s", cpuID, sl[1])
				core[coreID] = true
			case "vendor_id":
				if i.Data.CPU.Vendor == "" {
					i.Data.CPU.Vendor = sl[1]
				}
			case "model name":
				if i.Data.CPU.Model == "" {
					// CPU model, as reported by /proc/cpuinfo, can be a bit ugly. Clean up...
					model := reExtraSpace.ReplaceAllLiteralString(sl[1], " ")
					i.Data.CPU.Model = strings.Replace(model, "- ", "-", 1)
				}
			case "cache size":
				if i.Data.CPU.Cache == 0 {
					if m := reCacheSize.FindStringSubmatch(sl[1]); m != nil {
						if cache, err := strconv.ParseUint(m[1], 10, 64); err == nil {
							i.Data.CPU.Cache = uint(cache)
						}
					}
				}
			}
		}
	}
	if s.Err() != nil {
		return
	}

	i.Data.CPU.Cpus = uint(len(cpu))
	i.Data.CPU.Cores = uint(len(core))

	return
}
