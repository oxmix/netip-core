package collector

import (
	"sync"
	"time"
)

type CollectCore struct {
	Time     time.Time
	LoadAvg  []string
	CPUStats struct {
		Cores []int   `json:"cores"`
		Avg   float32 `json:"avg"`
	} `json:"cpuStats"`
	MemStats struct {
		MemTotal  int `json:"memTotal"`
		MemFree   int `json:"memFree"`
		Buffers   int `json:"buffers"`
		Cached    int `json:"cached"`
		Slab      int `json:"slab"`
		SwapTotal int `json:"swapTotal"`
		SwapFree  int `json:"swapFree"`
	} `json:"memStats"`
	IOStats map[string]IOStat `json:"ioStats"`
}

type IOStat struct {
	ReadIOPS       int `json:"readIOPS"`
	WriteIOPS      int `json:"writeIOPS"`
	DiscardIOPS    int `json:"discardIOPS"`
	ReadKbs        int `json:"readKbs"`
	WriteKbs       int `json:"writeKbs"`
	DiscardKbs     int `json:"discardKbs"`
	AwaitReadMs    int `json:"awaitReadMs"`
	AwaitWriteMs   int `json:"awaitWriteMs"`
	AwaitDiscardMs int `json:"awaitDiscardMs"`
	Utils          int `json:"utils"`
}

type WhoLogged struct {
	Time   time.Time `json:"time"`
	Device string    `json:"device"`
	User   string    `json:"user"`
	IP     string    `json:"ip"`
}

type Proc struct {
	PID     int    `json:"pid"`
	PPID    int    `json:"ppid"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Threads int    `json:"threads"`
	FDs     int    `json:"fds"`
}

type Collector struct {
	mu   sync.Mutex
	data CollectCore

	ChanWhoLogged chan *WhoLogged
	whoModTime    int64
	whoPointTime  time.Time

	ChanProcesses chan *Processes
}

func New() *Collector {
	c := new(Collector)

	go c.collectLoadAvg()
	go c.collectCPU()
	go c.collectMem()
	go c.collectIO()
	go c.collectWho()
	go c.collectProc()

	return c
}

func (c *Collector) Get() CollectCore {
	return c.data
}
