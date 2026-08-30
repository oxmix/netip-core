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
	Name          string `json:"name"`
	Version       string `json:"version"`
	SourceName    string `json:"sourceName"`
	SourceVersion string `json:"sourceVersion"`
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
		if force.Load() || time.Since(f.ModTime()) <= 3*time.Second {
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
				log.Printf("[collector] warn: no available packages parser for os: %q", osID)
				return
			}

			c.ChanPackages <- ipk
		}

		force.Store(false)
	}
}

func parseDebianLikeDpkg(path string) ([]Package, error) {
	ext, _ := loadExtendedStates("/apt-states")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pkgList []Package
	var current Package
	var installed bool

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "Package: "):
			current = Package{}
			installed = false
			current.Name = strings.TrimPrefix(line, "Package: ")

		case strings.HasPrefix(line, "Status: "):
			if strings.Contains(line, "install ok installed") {
				installed = true
			}

		case strings.HasPrefix(line, "Version: "):
			current.Version = strings.TrimPrefix(line, "Version: ")

		case strings.HasPrefix(line, "Source: "):

			src := strings.TrimPrefix(line, "Source: ")

			// samba (2:4.17.12+dfsg-0+deb12u3)
			if i := strings.Index(src, "("); i != -1 {
				current.SourceName = strings.TrimSpace(src[:i])
				current.SourceVersion = strings.TrimSuffix(src[i+1:], ")")
			} else {
				current.SourceName = src
			}

		case line == "":
			if installed {
				if len(ext) > 0 && !ext[current.Name] {
					continue
				}
				if current.SourceName != "" {
					if current.SourceVersion == "" {
						current.SourceVersion = current.Version
					}
					// wrapper packages where binary contains security revision
					if strings.HasPrefix(current.Version, current.SourceVersion+"+") {
						current.SourceVersion = current.Version
					}
				}
				pkgList = append(pkgList, current)
			}
		}
	}

	return pkgList, scanner.Err()
}

func loadExtendedStates(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pkgList := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	var pkg string

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Package: "):
			pkg = strings.TrimPrefix(line, "Package: ")
		case line == "":
			if pkg != "" {
				pkgList[pkg] = true
				pkg = ""
			}
		}
	}
	return pkgList, scanner.Err()
}
