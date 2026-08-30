package collector

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var (
	ioLoad = map[string]*IOStat{}
	ioPrev = map[string]*IOStat{}
)

func (c *Collector) collectIO() {
	for range time.Tick(time.Second) {
		c.disksStatsHandler("/proc/diskstats")
	}
}

func (c *Collector) disksStatsHandler(file string) {
	diskStats, err := os.Open(file)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(diskStats)

	defer c.mu.Unlock()
	c.mu.Lock()

	s := bufio.NewScanner(diskStats)
	for s.Scan() {
		var (
			devName                                              string
			major, minor, iosPgr, totTicks, rqTicks              int
			rdIos, rdMergesOrRdSec, rdSecOrWrIos, rdTicksOrWrSec int
			wrIos, wrMerges, wrSec, wrTicks                      int
			dcIos, dcMerges, dcSec, dcTicks                      int
			flIos, flTicks                                       int
		)

		_, err := fmt.Sscanf(s.Text(),
			"%d %d %s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
			&major, &minor, &devName,
			&rdIos, &rdMergesOrRdSec, &rdSecOrWrIos, &rdTicksOrWrSec,
			&wrIos, &wrMerges, &wrSec, &wrTicks, &iosPgr, &totTicks, &rqTicks,
			&dcIos, &dcMerges, &dcSec, &dcTicks,
			&flIos, &flTicks)

		if rdIos+wrIos+dcIos <= 0 {
			continue
		}

		_, err = os.Stat("/sys/block/" + devName)
		if err != nil {
			if os.IsNotExist(err) {
				if _, ok := ioPrev[devName]; ok {
					delete(ioPrev, devName)
					delete(ioLoad, devName)
				}
			}
			continue
		}

		eraseIOLoad := false
		if _, ok := ioPrev[devName]; !ok {
			eraseIOLoad = true
			ioPrev[devName] = &IOStat{}
			ioLoad[devName] = &IOStat{}
		}

		// IOps
		ioLoad[devName].ReadIOPS = rdIos - ioPrev[devName].ReadIOPS
		ioPrev[devName].ReadIOPS = rdIos

		ioLoad[devName].WriteIOPS = wrIos - ioPrev[devName].WriteIOPS
		ioPrev[devName].WriteIOPS = wrIos

		ioLoad[devName].DiscardIOPS = dcIos - ioPrev[devName].DiscardIOPS
		ioPrev[devName].DiscardIOPS = dcIos

		// speed
		rd := rdSecOrWrIos
		if rd > 0 {
			rd /= 2
		}
		ioLoad[devName].ReadKbs = rd - ioPrev[devName].ReadKbs
		ioPrev[devName].ReadKbs = rd

		wr := wrSec
		if wr > 0 {
			wr /= 2
		}
		ioLoad[devName].WriteKbs = wr - ioPrev[devName].WriteKbs
		ioPrev[devName].WriteKbs = wr

		dc := dcSec
		if dc > 0 {
			dc /= 2
		}
		ioLoad[devName].DiscardKbs = dc - ioPrev[devName].DiscardKbs
		ioPrev[devName].DiscardKbs = dc

		// await
		if ioLoad[devName].ReadIOPS > 0 {
			ioLoad[devName].AwaitReadMs =
				(rdTicksOrWrSec - ioPrev[devName].AwaitReadMs) / ioLoad[devName].ReadIOPS
		} else {
			ioLoad[devName].AwaitReadMs = 0
		}
		ioPrev[devName].AwaitReadMs = rdTicksOrWrSec

		if ioLoad[devName].WriteIOPS > 0 {
			ioLoad[devName].AwaitWriteMs =
				(wrTicks - ioPrev[devName].AwaitWriteMs) / ioLoad[devName].WriteIOPS
		} else {
			ioLoad[devName].AwaitWriteMs = 0
		}
		ioPrev[devName].AwaitWriteMs = wrTicks

		if ioLoad[devName].DiscardIOPS > 0 {
			ioLoad[devName].AwaitDiscardMs =
				(dcTicks - ioPrev[devName].AwaitDiscardMs) / ioLoad[devName].DiscardIOPS
		} else {
			ioLoad[devName].AwaitDiscardMs = 0
		}
		ioPrev[devName].AwaitDiscardMs = dcTicks

		// utilization
		tt := totTicks
		if tt > 0 {
			tt /= 10
		}
		ioLoad[devName].Utils = tt - ioPrev[devName].Utils
		ioPrev[devName].Utils = tt

		if eraseIOLoad {
			ioLoad[devName] = &IOStat{}
		}
	}

	ios := map[string]IOStat{}
	for devName, val := range ioLoad {
		ios[devName] = *val
	}
	c.data.IOStats = ios
}
