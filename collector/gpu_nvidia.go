package collector

import (
	"bufio"
	"encoding/xml"
	"log"
	"os/exec"
	"strings"
	"time"
)

var chanGpuNvidia = make(chan string, 1)

func (c *Collector) collectGpuNvidia() {
	cmd := exec.Command("nvidia-smi", "-q", "-x", "-l", "1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}
	go func() {
		for {
			c.parseGpuNvidia(<-chanGpuNvidia)
		}
	}()

	reader := bufio.NewReader(stdout)
	buff := strings.Builder{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		buff.WriteString(line)

		if strings.TrimSpace(line) == "</nvidia_smi_log>" {
			if len(chanGpuNvidia) == 0 {
				chanGpuNvidia <- buff.String()
			}
			buff.Reset()
		}

		// safe buffer flush over 3MB
		if buff.Len() > 3*1e6 {
			log.Println("warn: over buffer gpu stats, forced flush")
			buff.Reset()
		}
	}
}

type NvidiaSmiLog struct {
	XMLName       xml.Name `xml:"nvidia_smi_log"`
	Timestamp     string   `xml:"timestamp"`
	DriverVersion string   `xml:"driver_version"`
	CudaVersion   string   `xml:"cuda_version"`
	AttachedGpus  string   `xml:"attached_gpus"`
	Gpu           []struct {
		ID                  string `xml:"id,attr"`
		ProductName         string `xml:"product_name"`
		ProductBrand        string `xml:"product_brand"`
		ProductArchitecture string `xml:"product_architecture"`
		DisplayMode         string `xml:"display_mode"`
		DisplayActive       string `xml:"display_active"`
		PersistenceMode     string `xml:"persistence_mode"`
		MigMode             struct {
			CurrentMig string `xml:"current_mig"`
			PendingMig string `xml:"pending_mig"`
		} `xml:"mig_mode"`
		MigDevices               string `xml:"mig_devices"`
		AccountingMode           string `xml:"accounting_mode"`
		AccountingModeBufferSize string `xml:"accounting_mode_buffer_size"`
		DriverModel              struct {
			CurrentDm string `xml:"current_dm"`
			PendingDm string `xml:"pending_dm"`
		} `xml:"driver_model"`
		Serial           string `xml:"serial"`
		Uuid             string `xml:"uuid"`
		MinorNumber      string `xml:"minor_number"`
		VbiosVersion     string `xml:"vbios_version"`
		MultigpuBoard    string `xml:"multigpu_board"`
		BoardID          string `xml:"board_id"`
		BoardPartNumber  string `xml:"board_part_number"`
		GpuPartNumber    string `xml:"gpu_part_number"`
		GpuFruPartNumber string `xml:"gpu_fru_part_number"`
		GpuModuleID      string `xml:"gpu_module_id"`
		InforomVersion   struct {
			ImgVersion string `xml:"img_version"`
			OemObject  string `xml:"oem_object"`
			EccObject  string `xml:"ecc_object"`
			PwrObject  string `xml:"pwr_object"`
		} `xml:"inforom_version"`
		GpuOperationMode struct {
			CurrentGom string `xml:"current_gom"`
			PendingGom string `xml:"pending_gom"`
		} `xml:"gpu_operation_mode"`
		GspFirmwareVersion    string `xml:"gsp_firmware_version"`
		GpuVirtualizationMode struct {
			VirtualizationMode string `xml:"virtualization_mode"`
			HostVgpuMode       string `xml:"host_vgpu_mode"`
		} `xml:"gpu_virtualization_mode"`
		GpuResetStatus struct {
			ResetRequired            string `xml:"reset_required"`
			DrainAndResetRecommended string `xml:"drain_and_reset_recommended"`
		} `xml:"gpu_reset_status"`
		Ibmnpu struct {
			RelaxedOrderingMode string `xml:"relaxed_ordering_mode"`
		} `xml:"ibmnpu"`
		Pci struct {
			PciBus         string `xml:"pci_bus"`
			PciDevice      string `xml:"pci_device"`
			PciDomain      string `xml:"pci_domain"`
			PciDeviceID    string `xml:"pci_device_id"`
			PciBusID       string `xml:"pci_bus_id"`
			PciSubSystemID string `xml:"pci_sub_system_id"`
			PciGpuLinkInfo struct {
				PcieGen struct {
					MaxLinkGen           string `xml:"max_link_gen"`
					CurrentLinkGen       string `xml:"current_link_gen"`
					DeviceCurrentLinkGen string `xml:"device_current_link_gen"`
					MaxDeviceLinkGen     string `xml:"max_device_link_gen"`
					MaxHostLinkGen       string `xml:"max_host_link_gen"`
				} `xml:"pcie_gen"`
				LinkWidths struct {
					MaxLinkWidth     string `xml:"max_link_width"`
					CurrentLinkWidth string `xml:"current_link_width"`
				} `xml:"link_widths"`
			} `xml:"pci_gpu_link_info"`
			PciBridgeChip struct {
				BridgeChipType string `xml:"bridge_chip_type"`
				BridgeChipFw   string `xml:"bridge_chip_fw"`
			} `xml:"pci_bridge_chip"`
			ReplayCounter         string `xml:"replay_counter"`
			ReplayRolloverCounter string `xml:"replay_rollover_counter"`
			TxUtil                string `xml:"tx_util"`
			RxUtil                string `xml:"rx_util"`
			AtomicCapsInbound     string `xml:"atomic_caps_inbound"`
			AtomicCapsOutbound    string `xml:"atomic_caps_outbound"`
		} `xml:"pci"`
		FanSpeed              string `xml:"fan_speed"`
		PerformanceState      string `xml:"performance_state"`
		ClocksThrottleReasons struct {
			ClocksThrottleReasonGpuIdle                   string `xml:"clocks_throttle_reason_gpu_idle"`
			ClocksThrottleReasonApplicationsClocksSetting string `xml:"clocks_throttle_reason_applications_clocks_setting"`
			ClocksThrottleReasonSwPowerCap                string `xml:"clocks_throttle_reason_sw_power_cap"`
			ClocksThrottleReasonHwSlowdown                string `xml:"clocks_throttle_reason_hw_slowdown"`
			ClocksThrottleReasonHwThermalSlowdown         string `xml:"clocks_throttle_reason_hw_thermal_slowdown"`
			ClocksThrottleReasonHwPowerBrakeSlowdown      string `xml:"clocks_throttle_reason_hw_power_brake_slowdown"`
			ClocksThrottleReasonSyncBoost                 string `xml:"clocks_throttle_reason_sync_boost"`
			ClocksThrottleReasonSwThermalSlowdown         string `xml:"clocks_throttle_reason_sw_thermal_slowdown"`
			ClocksThrottleReasonDisplayClocksSetting      string `xml:"clocks_throttle_reason_display_clocks_setting"`
		} `xml:"clocks_throttle_reasons"`
		FbMemoryUsage struct {
			Total    string `xml:"total"`
			Reserved string `xml:"reserved"`
			Used     string `xml:"used"`
			Free     string `xml:"free"`
		} `xml:"fb_memory_usage"`
		Bar1MemoryUsage struct {
			Total string `xml:"total"`
			Used  string `xml:"used"`
			Free  string `xml:"free"`
		} `xml:"bar1_memory_usage"`
		ComputeMode string `xml:"compute_mode"`
		Utilization struct {
			GpuUtil     string `xml:"gpu_util"`
			MemoryUtil  string `xml:"memory_util"`
			EncoderUtil string `xml:"encoder_util"`
			DecoderUtil string `xml:"decoder_util"`
		} `xml:"utilization"`
		EncoderStats struct {
			SessionCount   string `xml:"session_count"`
			AverageFps     string `xml:"average_fps"`
			AverageLatency string `xml:"average_latency"`
		} `xml:"encoder_stats"`
		FbcStats struct {
			SessionCount   string `xml:"session_count"`
			AverageFps     string `xml:"average_fps"`
			AverageLatency string `xml:"average_latency"`
		} `xml:"fbc_stats"`
		EccMode struct {
			CurrentEcc string `xml:"current_ecc"`
			PendingEcc string `xml:"pending_ecc"`
		} `xml:"ecc_mode"`
		EccErrors struct {
			Volatile struct {
				SingleBit struct {
					DeviceMemory  string `xml:"device_memory"`
					RegisterFile  string `xml:"register_file"`
					L1Cache       string `xml:"l1_cache"`
					L2Cache       string `xml:"l2_cache"`
					TextureMemory string `xml:"texture_memory"`
					TextureShm    string `xml:"texture_shm"`
					Cbu           string `xml:"cbu"`
					Total         string `xml:"total"`
				} `xml:"single_bit"`
				DoubleBit struct {
					DeviceMemory  string `xml:"device_memory"`
					RegisterFile  string `xml:"register_file"`
					L1Cache       string `xml:"l1_cache"`
					L2Cache       string `xml:"l2_cache"`
					TextureMemory string `xml:"texture_memory"`
					TextureShm    string `xml:"texture_shm"`
					Cbu           string `xml:"cbu"`
					Total         string `xml:"total"`
				} `xml:"double_bit"`
			} `xml:"volatile"`
			Aggregate struct {
				SingleBit struct {
					DeviceMemory  string `xml:"device_memory"`
					RegisterFile  string `xml:"register_file"`
					L1Cache       string `xml:"l1_cache"`
					L2Cache       string `xml:"l2_cache"`
					TextureMemory string `xml:"texture_memory"`
					TextureShm    string `xml:"texture_shm"`
					Cbu           string `xml:"cbu"`
					Total         string `xml:"total"`
				} `xml:"single_bit"`
				DoubleBit struct {
					DeviceMemory  string `xml:"device_memory"`
					RegisterFile  string `xml:"register_file"`
					L1Cache       string `xml:"l1_cache"`
					L2Cache       string `xml:"l2_cache"`
					TextureMemory string `xml:"texture_memory"`
					TextureShm    string `xml:"texture_shm"`
					Cbu           string `xml:"cbu"`
					Total         string `xml:"total"`
				} `xml:"double_bit"`
			} `xml:"aggregate"`
		} `xml:"ecc_errors"`
		RetiredPages struct {
			MultipleSingleBitRetirement struct {
				RetiredCount    string `xml:"retired_count"`
				RetiredPagelist string `xml:"retired_pagelist"`
			} `xml:"multiple_single_bit_retirement"`
			DoubleBitRetirement struct {
				RetiredCount    string `xml:"retired_count"`
				RetiredPagelist string `xml:"retired_pagelist"`
			} `xml:"double_bit_retirement"`
			PendingBlacklist  string `xml:"pending_blacklist"`
			PendingRetirement string `xml:"pending_retirement"`
		} `xml:"retired_pages"`
		RemappedRows string `xml:"remapped_rows"`
		Temperature  struct {
			GpuTemp                string `xml:"gpu_temp"`
			GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
			GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
			GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
			GpuTargetTemperature   string `xml:"gpu_target_temperature"`
			MemoryTemp             string `xml:"memory_temp"`
			GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
		} `xml:"temperature"`
		SupportedGpuTargetTemp struct {
			GpuTargetTempMin string `xml:"gpu_target_temp_min"`
			GpuTargetTempMax string `xml:"gpu_target_temp_max"`
		} `xml:"supported_gpu_target_temp"`
		PowerReadings struct {
			PowerState         string `xml:"power_state"`
			PowerManagement    string `xml:"power_management"`
			PowerDraw          string `xml:"power_draw"`
			PowerLimit         string `xml:"power_limit"`
			DefaultPowerLimit  string `xml:"default_power_limit"`
			EnforcedPowerLimit string `xml:"enforced_power_limit"`
			MinPowerLimit      string `xml:"min_power_limit"`
			MaxPowerLimit      string `xml:"max_power_limit"`
		} `xml:"power_readings"`
		GpuPowerReadings struct {
			PowerState          string `xml:"power_state"`
			PowerDraw           string `xml:"power_draw"`
			CurrentPowerLimit   string `xml:"current_power_limit"`
			RequestedPowerLimit string `xml:"requested_power_limit"`
			DefaultPowerLimit   string `xml:"default_power_limit"`
			MinPowerLimit       string `xml:"min_power_limit"`
			MaxPowerLimit       string `xml:"max_power_limit"`
		} `xml:"gpu_power_readings"`
		Clocks struct {
			GraphicsClock string `xml:"graphics_clock"`
			SmClock       string `xml:"sm_clock"`
			MemClock      string `xml:"mem_clock"`
			VideoClock    string `xml:"video_clock"`
		} `xml:"clocks"`
		ApplicationsClocks struct {
			GraphicsClock string `xml:"graphics_clock"`
			MemClock      string `xml:"mem_clock"`
		} `xml:"applications_clocks"`
		DefaultApplicationsClocks struct {
			GraphicsClock string `xml:"graphics_clock"`
			MemClock      string `xml:"mem_clock"`
		} `xml:"default_applications_clocks"`
		DeferredClocks struct {
			MemClock string `xml:"mem_clock"`
		} `xml:"deferred_clocks"`
		MaxClocks struct {
			GraphicsClock string `xml:"graphics_clock"`
			SmClock       string `xml:"sm_clock"`
			MemClock      string `xml:"mem_clock"`
			VideoClock    string `xml:"video_clock"`
		} `xml:"max_clocks"`
		MaxCustomerBoostClocks struct {
			GraphicsClock string `xml:"graphics_clock"`
		} `xml:"max_customer_boost_clocks"`
		ClockPolicy struct {
			AutoBoost        string `xml:"auto_boost"`
			AutoBoostDefault string `xml:"auto_boost_default"`
		} `xml:"clock_policy"`
		Voltage struct {
			GraphicsVolt string `xml:"graphics_volt"`
		} `xml:"voltage"`
		Fabric struct {
			State  string `xml:"state"`
			Status string `xml:"status"`
		} `xml:"fabric"`
		SupportedClocks struct {
			SupportedMemClock []struct {
				Value                  string   `xml:"value"`
				SupportedGraphicsClock []string `xml:"supported_graphics_clock"`
			} `xml:"supported_mem_clock"`
		} `xml:"supported_clocks"`
		Processes struct {
			ProcessInfo []struct {
				GpuInstanceID     string `xml:"gpu_instance_id"`
				ComputeInstanceID string `xml:"compute_instance_id"`
				Pid               string `xml:"pid"`
				Type              string `xml:"type"`
				ProcessName       string `xml:"process_name"`
				UsedMemory        string `xml:"used_memory"`
			} `xml:"process_info"`
		} `xml:"processes"`
		AccountedProcesses string `xml:"accounted_processes"`
	} `xml:"gpu"`
}

func (c *Collector) parseGpuNvidia(data string) []SmiNvidia {
	smi := &NvidiaSmiLog{}
	err := xml.Unmarshal([]byte(data), smi)
	if err != nil {
		log.Println(err)
		return nil
	}

	var stats []SmiNvidia
	for _, g := range smi.Gpu {
		power := g.PowerReadings.PowerDraw
		powerLimit := g.PowerReadings.PowerLimit
		if power == "" && powerLimit == "" {
			// changing struct since NVIDIA-SMI 550.54.14 / Driver Version: 550.54.14
			power = g.GpuPowerReadings.PowerDraw
			powerLimit = g.GpuPowerReadings.CurrentPowerLimit
		}

		stats = append(stats, SmiNvidia{
			Name:   g.ProductName,
			Driver: smi.DriverVersion,
			PciDev: g.Pci.PciDevice,
			PciBus: g.Pci.PciBusID,
			PciLink: g.Pci.PciGpuLinkInfo.LinkWidths.CurrentLinkWidth +
				" " + g.Pci.PciGpuLinkInfo.PcieGen.DeviceCurrentLinkGen + ".0",
			PciTx:      g.Pci.TxUtil,
			PciRx:      g.Pci.RxUtil,
			Fan:        g.FanSpeed,
			MemUse:     g.FbMemoryUsage.Used,
			MemFree:    g.FbMemoryUsage.Free,
			UtilGpu:    g.Utilization.GpuUtil,
			UtilMem:    g.Utilization.MemoryUtil,
			UtilEnc:    g.Utilization.EncoderUtil,
			UtilDec:    g.Utilization.DecoderUtil,
			TmpGpu:     g.Temperature.GpuTemp,
			TmpMem:     g.Temperature.MemoryTemp,
			Power:      power,
			PowerLimit: powerLimit,
			ClockGra:   g.Clocks.GraphicsClock,
			ClockSm:    g.Clocks.SmClock,
			ClockMem:   g.Clocks.MemClock,
			ClockVideo: g.Clocks.VideoClock,
		})
	}

	defer c.mu.Unlock()
	c.mu.Lock()
	c.data.Time = time.Now().UTC()
	c.data.GPUStats.Nvidia = stats

	return stats
}
