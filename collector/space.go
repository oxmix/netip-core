package collector

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func (c *Collector) collectSpace() {
	var volumes []string
	if _, err := os.Stat("/_external"); !os.IsNotExist(err) {
		err := filepath.Walk("/_external/", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() != ".netip-device" {
				return nil
			}
			volumes = append(volumes, path)
			return nil
		})
		if err != nil {
			log.Println("collect space err: ", err)
		}
	}

	for {
		c.data.Time = time.Now().UTC()
		c.data.SpaceStats = map[string]SpaceStatFS{}
		c.handlerStatFS("/")
		for _, path := range volumes {
			c.handlerStatFS(path)
		}
		time.Sleep(time.Second)
	}
}

func (c *Collector) handlerStatFS(path string) {
	var sf syscall.Statfs_t
	err := syscall.Statfs(path, &sf)
	if err != nil {
		log.Println("stat fs path: "+path+" err: ", err)
		return
	}
	r := SpaceStatFS{
		ReservedBlocks: sf.Bfree - sf.Bavail,
		BlockSize:      uint64(sf.Bsize),
	}
	r.Total = r.BlockSize * (sf.Blocks - r.ReservedBlocks)
	r.Free = r.BlockSize * sf.Bavail

	defer c.mu.Unlock()
	c.mu.Lock()
	realPath := strings.Replace(path, "/_external", "", 1)
	c.data.SpaceStats[filepath.Dir(realPath)] = r
}
