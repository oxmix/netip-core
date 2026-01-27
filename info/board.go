package info

import (
	"os"
	"strings"
)

func (i *Info) fillBoardInfo() {
	name, err := os.ReadFile("/sys/class/dmi/id/board_name")
	if err == nil {
		i.Data.Board.Name = strings.TrimSpace(string(name))
	}

	vendor, err := os.ReadFile("/sys/class/dmi/id/board_vendor")
	if err == nil {
		i.Data.Board.Vendor = strings.TrimSpace(string(vendor))
	}

	biosVersion, err := os.ReadFile("/sys/class/dmi/id/bios_version")
	if err == nil {
		i.Data.Board.BiosVersion = strings.TrimSpace(string(biosVersion))

		biosDate, err := os.ReadFile("/sys/class/dmi/id/bios_date")
		if err == nil {
			i.Data.Board.BiosVersion += " " + string(biosDate)
		}
	}
}
