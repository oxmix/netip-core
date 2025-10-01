package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *Collector) collectGpuAmd() {
	var cards []uint8
	for i := 0; i < 10; i++ {
		file := fmt.Sprintf("/sys/class/drm/card%d/device/gpu_metrics", i)
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		if len(data) < 128 {
			log.Printf("warn: gpu_metrics too small: %d bytes", len(data))
			continue
		}
		cards = append(cards, uint8(i))
	}

	if len(cards) == 0 {
		return
	}

	for {
		sas := make([]GpuAmd, 0, len(cards))
		for _, card := range cards {
			data, _ := os.ReadFile(fmt.Sprintf("/sys/class/drm/card%d/device/gpu_metrics", card))
			if pga, err := c.parseGpuAmd(card, data); err == nil {
				sas = append(sas, pga)
			}
		}
		c.writeAmdStats(sas)
		time.Sleep(time.Second)
	}
}

type AMDGpuMetricsV22 struct {
	StructureSize   uint16
	FormatRevision  uint8
	ContentRevision uint8

	TemperatureGfx  uint16
	TemperatureSoc  uint16
	TemperatureCore [8]uint16
	TemperatureL3   [2]uint16

	AverageGfxActivity uint16
	AverageMmActivity  uint16

	SystemClockCounter uint64

	AverageSocketPower uint16
	AverageCpuPower    uint16
	AverageSocPower    uint16
	AverageGfxPower    uint16
	AverageCorePower   [8]uint16

	AverageGfxClk uint16
	AverageSocClk uint16
	AverageUClk   uint16
	AverageFClk   uint16
	AverageVClk   uint16
	AverageDClk   uint16

	CurrentGfxClk  uint16
	CurrentSocClk  uint16
	CurrentUClk    uint16
	CurrentFClk    uint16
	CurrentVClk    uint16
	CurrentDClk    uint16
	CurrentCoreClk [8]uint16
	CurrentL3Clk   [2]uint16

	ThrottleStatus uint32
	FanPwm         uint16
	Padding        [3]uint16

	IndependentThrottleStatus uint64
}

func (c *Collector) parseGpuAmd(card uint8, data []byte) (GpuAmd, error) {
	var m AMDGpuMetricsV22
	if err := binary.Read(
		bytes.NewReader(data),
		binary.LittleEndian,
		&m,
	); err != nil {
		return GpuAmd{}, err
	}

	vendor, _ := c.readHexStringFromFile(fmt.Sprintf("/sys/class/drm/card%d/device/vendor", card))
	device, _ := c.readHexStringFromFile(fmt.Sprintf("/sys/class/drm/card%d/device/device", card))
	memUsed, _ := c.readUint64FromFile(fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_used", card))
	memTotal, _ := c.readUint64FromFile(fmt.Sprintf("/sys/class/drm/card%d/device/mem_info_vram_total", card))

	return GpuAmd{
		VerFormat:  m.FormatRevision,
		VerContent: m.ContentRevision,
		Card:       card,
		Vendor:     vendor,
		Device:     device,
		MemUse:     memUsed,
		MemFree:    memTotal - memUsed,
		UtilGpu:    fmt.Sprintf("%.2f", float64(m.AverageGfxActivity)/100),
		UtilMedia:  fmt.Sprintf("%.2f", float64(m.AverageMmActivity)/100),
		TmpGpu:     fmt.Sprintf("%.2f", float64(m.TemperatureSoc)/100),
		Power:      fmt.Sprintf("%.2f", float64(m.AverageSocPower)/100),
		ClockSoc:   int(m.CurrentSocClk),
	}, nil
}

func (c *Collector) writeAmdStats(data []GpuAmd) {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.data.Time = time.Now().UTC()
	c.data.GPUStats.Amd = data
}

func (c *Collector) readHexStringFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %v", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func (c *Collector) readUint64FromFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %v", path, err)
	}
	s := strings.TrimSpace(string(data))
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("gpuAmd: failed to parse %s: %v", path, err)
	}
	return val, nil
}
