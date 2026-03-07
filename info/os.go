package info

import (
	"bufio"
	"os"
	"strings"
)

func (i *Info) fillOSInfo() {
	file, err := os.Open("/etc/host-os-release")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			i.Data.OS.ID = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		}
		if strings.HasPrefix(line, "NAME=") {
			i.Data.OS.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			i.Data.OS.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_CODENAME=") {
			i.Data.OS.Codename = strings.Trim(strings.TrimPrefix(line, "VERSION_CODENAME="), `"`)
		}
	}
}
