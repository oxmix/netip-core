package collector

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const pathHwm = "/sys/class/hwmon"

type TempStats struct {
	Label    string  `json:"label,omitempty"`
	Temp     float64 `json:"temp,omitempty"`
	TempMax  float64 `json:"max,omitempty"`
	TempCrit float64 `json:"crit,omitempty"`
}

type prepareHWM struct {
	TempPath string
	Label    string
	TempMax  float64
	TempCrit float64
}

func parseSysfsMilliDegrees(raw []byte) (float64, error) {
	s := strings.TrimSpace(string(raw))
	if s == "" {
		return 0, fmt.Errorf("empty temperature value")
	}
	milli, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return milli / 1000.0, nil
}

func (c *Collector) collectTemperature() {
	list, err := os.ReadDir(pathHwm)
	if err != nil || len(list) == 0 {
		return
	}

	var HWMs []prepareHWM
	for _, entry := range list {
		id := strings.TrimPrefix(entry.Name(), "hwmon")
		if id == "" {
			log.Println("[collector] got zero hwm id")
			continue
		}

		hwDir := pathHwm + "/" + entry.Name()

		files, err := filepath.Glob(hwDir + "/temp*_input")
		if err != nil || len(files) == 0 {
			continue
		}

		for _, file := range files {
			n := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), "temp"), "_input")
			if n == "" {
				continue
			}
			hwm := prepareHWM{
				TempPath: file,
			}
			tLabel, _ := os.ReadFile(hwDir + "/temp" + n + "_label")
			tMax, _ := os.ReadFile(hwDir + "/temp" + n + "_max")
			if len(tMax) > 0 {
				if f, err := parseSysfsMilliDegrees(tMax); err == nil {
					hwm.TempMax = f
				}
			}
			tCrit, _ := os.ReadFile(hwDir + "/temp" + n + "_crit")
			if len(tCrit) > 0 {
				if f, err := parseSysfsMilliDegrees(tCrit); err == nil {
					hwm.TempCrit = f
				}
			}

			label := bytes.Buffer{}
			label.WriteString("hwm")
			label.WriteString(id)
			label.WriteString(" ")
			name, _ := os.ReadFile(hwDir + "/name")
			if len(name) != 0 {
				nm := strings.TrimSuffix(strings.TrimSpace(strings.ToLower(string(name))), "temp")
				if nm == "" {
					nm = strings.TrimSpace(strings.ToLower(string(name)))
				}
				label.WriteString(nm)
				label.WriteString(" ")
			}
			if len(tLabel) > 0 {
				label.WriteString(strings.TrimSpace(strings.ToLower(string(tLabel))))
			}
			hwm.Label = strings.TrimSpace(label.String())

			HWMs = append(HWMs, hwm)
		}
	}

	if len(HWMs) > 0 {
		c.tempStatHandler(HWMs)
	}
}

func (c *Collector) tempStatHandler(HWMs []prepareHWM) {
	for range time.Tick(time.Second) {
		stats := make([]TempStats, 0, len(HWMs))
		for _, hwm := range HWMs {
			tInp, err := os.ReadFile(hwm.TempPath)
			if err != nil {
				continue
			}
			temp, err := parseSysfsMilliDegrees(tInp)
			if err == nil {
				stats = append(stats, TempStats{
					Label:    hwm.Label,
					Temp:     temp,
					TempMax:  hwm.TempMax,
					TempCrit: hwm.TempCrit,
				})
			}
		}
		c.mu.Lock()
		c.data.TempStats = stats
		c.mu.Unlock()
	}
}
