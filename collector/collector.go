package collector

import (
	"log"
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
	IOStats    map[string]IOStat      `json:"ioStats"`
	SpaceStats map[string]SpaceStatFS `json:"spaceStats"`
	GPUStats   GPUStats               `json:"gpuStats"`
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
	Session string    `json:"session"`
	Time    time.Time `json:"time"`
	Device  string    `json:"device"`
	User    string    `json:"user"`
	IP      string    `json:"ip"`
}

type Proc struct {
	PID     int    `json:"pid"`
	PPID    int    `json:"ppid"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Threads int    `json:"threads"`
	FDs     int    `json:"fds"`
}

type SpaceStatFS struct {
	ReservedBlocks uint64 `json:"reservedBlocks"`
	BlockSize      uint64 `json:"blockSize"`
	Total          uint64 `json:"total"`
	Free           uint64 `json:"free"`
}

type SmartDisk struct {
	Model       string `json:"modelFamily"`
	Serial      string `json:"serialNumber"`
	Capacity    string `json:"capacity"`
	Health      string `json:"health"`
	Used        string `json:"used"`
	Working     int64  `json:"working"`
	Temperature string `json:"temperature"`
	Full        string `json:"full"`
	Error       string `json:"error"`
}

type RaidMD struct {
	Disks   []string  `json:"disks"`
	Proc    RaidProc  `json:"proc"`
	ProcOut string    `json:"procOut"`
	Adm     RaidMDAdm `json:"adm"`
	AdmOut  string    `json:"admOut"`
}

type RaidProc struct {
	State    string  `json:"state"`
	Progress float64 `json:"progress"`
	Left     float64 `json:"left"`
	Speed    int     `json:"speed"`
	ParseErr string  `json:"parseErr,omitempty"`
}

type RaidMDAdm struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	Level     string `json:"level"`
	Capacity  string `json:"capacity"`
	CreatedAt string `json:"createdAt"`
	Active    string `json:"active"`
	Working   string `json:"working"`
	Failed    string `json:"failed"`
	Spare     string `json:"spare"`
}

type DisksInfo struct {
	Version int                   `json:"version"`
	Time    time.Time             `json:"time"`
	Smarts  map[string]*SmartDisk `json:"smarts"`
	Raids   map[string]*RaidMD    `json:"raids"`
}

type GPUStats struct {
	Nvidia []SmiNvidia `json:"nvidia,omitempty"`
	Amd    []GpuAmd    `json:"amd,omitempty"`
}

type SmiNvidia struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	PciDev     string `json:"pciDev"`
	PciBus     string `json:"pciBus"`
	PciLink    string `json:"pciLink"`
	PciTx      string `json:"pciTx"`
	PciRx      string `json:"pciRx"`
	Fan        string `json:"fan"`
	MemUse     string `json:"memUse"`
	MemFree    string `json:"memFree"`
	UtilGpu    string `json:"utilGpu"`
	UtilMem    string `json:"utilMem"`
	UtilEnc    string `json:"utilEnc"`
	UtilDec    string `json:"utilDec"`
	TmpGpu     string `json:"tmpGpu"`
	TmpMem     string `json:"tmpMem"`
	Power      string `json:"power"`
	PowerLimit string `json:"powerLimit"`
	ClockGra   string `json:"clockSh"`
	ClockSm    string `json:"clockSm"`
	ClockMem   string `json:"clockMem"`
	ClockVideo string `json:"clockVideo"`
}

type GpuAmd struct {
	VerFormat  uint8  `json:"verFormat"`
	VerContent uint8  `json:"verContent"`
	Card       uint8  `json:"card"`
	Vendor     string `json:"vendor"`
	Device     string `json:"device"`
	MemUse     uint64 `json:"memUse"`
	MemFree    uint64 `json:"memFree"`
	UtilGpu    string `json:"utilGpu"`
	UtilMedia  string `json:"utilMedia"`
	TmpGpu     string `json:"tmpGpu"`
	Power      string `json:"power"`
	ClockSoc   int    `json:"clockSoc"`
}

type Collector struct {
	mu       sync.RWMutex
	data     CollectCore
	ChanCore chan *CollectCore

	ChanWhoLogged chan *WhoLogged
	whoSessionMu  sync.Mutex
	whoSessions   map[string]*WhoLogged

	ChanProcesses chan *Processes
	ChanDisksInfo chan *DisksInfo
}

func New() *Collector {
	c := &Collector{
		mu:       sync.RWMutex{},
		data:     CollectCore{},
		ChanCore: make(chan *CollectCore, 1),

		ChanWhoLogged: make(chan *WhoLogged, 128),

		ChanProcesses: make(chan *Processes, 1),
		ChanDisksInfo: make(chan *DisksInfo, 1),
	}

	go c.senderCore()
	go c.collectLoadAvg()
	go c.collectCPU()
	go c.collectMem()
	go c.collectIO()
	go c.collectWho()
	go c.collectProc()
	go c.collectSpace()
	go c.collectDisks()
	go c.collectGpuNvidia()
	go c.collectGpuAmd()

	return c
}

func (c *Collector) senderCore() {
	for range time.Tick(time.Second) {
		c.mu.RLock()
		snapshot := c.data // shallow copy of struct
		c.mu.RUnlock()

		if len(c.ChanCore) > 0 {
			log.Println("[collector] notice: core chan is throttling")
		}

		// use this, for uniq address: snap := snapshot;
		c.ChanCore <- &snapshot
	}
}
