package collector

import (
	"log"
	u "netip-core/utils"
	"os"
	"strings"
	"time"
)

func (c *Collector) collectWho() {
	c.whoPointTime = time.Now().UTC()
	c.ChanWhoLogged = make(chan *WhoLogged, 255)

	file := "/var/run/utmp"
	for {
		stat, err := os.Stat(file)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		modTime := stat.ModTime().Unix()
		if c.whoModTime == modTime {
			time.Sleep(3 * time.Second)
			continue
		}
		c.whoModTime = modTime

		mp, err := u.ReadUtmp(file)
		if err != nil {
			log.Println("parse utmp err: ", err)
			continue
		}

		if len(mp) == 0 {
			continue
		}

		c.whoWrite(mp)

		c.whoPointTime = mp[len(mp)-1].Time
	}
}

func (c *Collector) whoWrite(mp []*u.UtmpSmart) {
	for _, a := range mp {
		// already skip
		if !a.Time.After(c.whoPointTime) {
			continue
		}

		// utmpdump /var/run/utmp
		//#define INIT_PROCESS    5
		//#define LOGIN_PROCESS   6
		//#define USER_PROCESS    7
		//#define DEAD_PROCESS    8
		if a.Type != 7 {
			continue
		}

		// pts type logged from local
		if strings.HasPrefix(a.Device, "pts/") && a.Addr == "0.0.0.0" {
			continue
		}

		c.ChanWhoLogged <- &WhoLogged{
			Time:   a.Time,
			Device: a.Device,
			User:   a.User,
			IP:     a.Addr,
		}
	}
}
