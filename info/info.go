package info

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	Uptime int64 `json:"uptime"`
	Data   struct {
		CPU struct {
			Vendor  string `json:"vendor"`
			Model   string `json:"model"`
			Speed   uint   `json:"speed"`   // clock rate in MHz
			Cache   uint   `json:"cache"`   // cache size in KB
			Cpus    uint   `json:"cpus"`    // physical CPUs
			Cores   uint   `json:"cores"`   // physical CPU cores
			Threads uint   `json:"threads"` // logical (HT) CPU cores
		} `json:"cpu"`
		Board struct {
			Name        string `json:"name"`
			Vendor      string `json:"vendor"`
			BiosVersion string `json:"biosVersion"`
		} `json:"board"`
		Kernel struct {
			Architecture string `json:"architecture"`
			OSType       string `json:"osType"`
			OSRelease    string `json:"osRelease"`
			OSVersion    string `json:"osVersion"`
		} `json:"kernel"`
	} `json:"data"`
}

func Get() (i *Info) {
	i = new(Info)

	i.fillUptime()
	i.fillCPUInfo()
	i.fillBoardInfo()
	i.fillKernelInfo()

	return
}

func (i *Info) fillUptime() {
	timeG, _ := os.ReadFile("/proc/uptime")
	timeF, _ := strconv.ParseFloat(strings.Split(string(timeG), " ")[0], 64)
	i.Uptime = time.Now().Unix() - int64(timeF)
}
