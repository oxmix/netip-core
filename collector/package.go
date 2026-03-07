package collector

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Collector) collectPackage(osID string) {
	force := atomic.Bool{}
	force.Store(true)
	go func() {
		for range time.Tick(12 * time.Hour) {
			force.Store(true)
		}
	}()

	for range time.Tick(3 * time.Second) {
		file := "/dpkg-status"
		f, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			log.Printf("[collector] dpkg status file read err: %v", err)
			return
		}

		var ipk []Package
		if force.Load() || time.Now().Sub(f.ModTime()).Seconds() <= 3 {
			switch osID {
			case "debian", "kali", "raspbian", "parrot", "devuan", "astra":
				ipk, err = parseDebianLikeDpkg("/dpkg-status")
				if err != nil {
					log.Printf("[collector] parseDebianLikeDpkg, osID: %s err: %v", osID, err)
				}
			case "ubuntu", "linuxmint", "pop", "zorin", "elementary", "neon":
				ipk, err = parseDebianLikeDpkg("/dpkg-status")
				if err != nil {
					log.Printf("[collector] parseDebianLikeDpkg, osID: %s err: %v", osID, err)
				}
			default:
				log.Println("[collector] warn: no available packages parser for os:", osID)
				return
			}

			c.ChanPackages <- ipk
		}

		force.Store(false)
	}
}

func parseDebianLikeDpkg(path string) ([]Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var PKGs []Package
	var current Package
	var installed bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Package: ") {
			current = Package{}
			installed = false
			current.Name = strings.TrimPrefix(line, "Package: ")
		}

		if strings.HasPrefix(line, "Status: ") {
			if strings.Contains(line, "install ok installed") {
				installed = true
			}
		}

		if strings.HasPrefix(line, "Version: ") {
			current.Version = strings.TrimPrefix(line, "Version: ")
		}

		// end block pack
		if line == "" {
			if installed && current.Name != "" {
				PKGs = append(PKGs, current)
				current = Package{}
			}
		}
	}

	return PKGs, scanner.Err()
}
