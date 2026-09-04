package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"netip-core/libs/logger"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const diskCmdTimeout = 30 * time.Second

func (c *Collector) collectDisks() {
	c.collectDisksOnce()
	for range time.Tick(15 * time.Minute) {
		c.collectDisksOnce()
	}
}

func (c *Collector) listDisks() []string {
	ctx, cancel := context.WithTimeout(context.Background(), diskCmdTimeout)
	defer cancel()

	var disks []string
	output, err := exec.CommandContext(ctx, "sh", "-c",
		"lsblk -d -n -o NAME,RO | awk '/0$/ {print $1}'").Output()
	if err != nil {
		log.Println("[collector] lsblk err:", err)
		return disks
	}
	for _, d := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if d == "" {
			continue
		}
		disks = append(disks, d)
	}
	return disks
}

func (c *Collector) collectDisksOnce() {
	disks := c.listDisks()

	di := &DisksInfo{
		Version: 2,
		Time:    time.Now().UTC(),
		Smarts:  map[string]*SmartDisk{},
		Raids:   map[string]*RaidMD{},
		Zfs:     []RaidZFS{},
	}

	// smarts disks
	for _, dev := range disks {
		ctx, cancel := context.WithTimeout(context.Background(), diskCmdTimeout)
		info, err := exec.CommandContext(ctx, "sh", "-c", "smartctl --all /dev/"+dev).CombinedOutput()
		cancel()
		if logger.IsDebugMode() {
			logger.Debugf("[collector] smartctl --all /dev/%s err: %s out: %s", dev, err, info)
		}
		if err != nil {
			continue
		}
		di.Smarts[dev] = c.parseSmart(string(info))
	}

	// raids md
	mdStat, err := os.ReadFile("/proc/mdstat")
	if err == nil {
		di.Raids = c.parseMdStat(string(mdStat))
	}

	for md, adm := range di.Raids {
		ctx, cancel := context.WithTimeout(context.Background(), diskCmdTimeout)
		mdAdm, err := exec.CommandContext(ctx, "sh", "-c", "mdadm -D /dev/"+md).Output()
		cancel()
		if err == nil {
			if logger.IsDebugMode() {
				logger.Debugf("[collector] mdadm raw: %s", mdAdm)
			}
			adm.Adm = c.parseMdAdm(string(mdAdm))
		}
	}

	// raids zfs
	zfsCtx, zfsCancel := context.WithTimeout(context.Background(), diskCmdTimeout)
	zfsListJson, err1 := exec.CommandContext(zfsCtx, "sh", "-c", "zpool list -vPpj").Output()
	zfsStatusJson, err2 := exec.CommandContext(zfsCtx, "sh", "-c", "zpool status -jP").Output()
	zfsCancel()
	if err1 == nil && err2 == nil {
		di.Zfs, err = c.parseZfs(zfsListJson, zfsStatusJson)
		if err != nil {
			log.Println("[collector] zfs parse err:", err)
		}
	}

	c.ChanDisksInfo <- di
}

var reSmartModel = regexp.MustCompile(`Model Family:(.*?)\nDevice Model:(.*?)\n`)
var reSmartModel2 = regexp.MustCompile(`(?m)Model Number:(.*?)$`)
var reSmartModel3 = regexp.MustCompile(`(?m)Device Model:(.*?)$`)
var reSmartSerial = regexp.MustCompile(`(?m)Serial Number:(.*?)$`)

var reSmartCapacity = regexp.MustCompile(`(?m)User Capacity:.*\[(.*?)]$`)
var reSmartCapacity2 = regexp.MustCompile(`(?m)Total NVM Capacity:.*\[(.*?)]$`)
var reSmartCapacity3 = regexp.MustCompile(`(?m)Namespace 1 Size/Capacity:.*\[(.*?)]$`)

var reSmartUsed = regexp.MustCompile(`(?m)Percentage Used:(.*%)$`)
var reSmartUsedAttr231 = regexp.MustCompile(`(?m)^\s*231\s+[\w_]+\s+\S+\s+(\d+)\s+`)
var reSmartUsedAttr233 = regexp.MustCompile(`(?m)^\s*233\s+[\w_]+\s+\S+\s+(\d+)\s+`)

var reSmartHealth = regexp.MustCompile(`(?m)SMART overall-health self-assessment test result: (.*?)$`)
var reSmartHealth2 = regexp.MustCompile(`(?m)SMART Health Status: (.*?)$`)

var reSmartWorking = regexp.MustCompile(`(?m)Power_On_Hours.*?-.*?(.*?)(?:\(|$)`)
var reSmartWorking2 = regexp.MustCompile(`(?m)Power On Hours:.*?(.*?)$`)

var reSmartTemperature = regexp.MustCompile(`(?m)(?:Airflow_Temperature_Cel|Temperature_Case).*?-.*?(.*?)$`)
var reSmartTemperature2 = regexp.MustCompile(`(?m)Temperature_Celsius.*?-.*?(.*?)(?:\(|$)$`)
var reSmartTemperature3 = regexp.MustCompile(`Temperature Sensor.*?:.*?(.*)Cel`)
var reSmartTemperature4 = regexp.MustCompile(`Temperature.*?:.*?(.*)Celsius`)

func (c *Collector) parseSmart(data string) *SmartDisk {
	sd := &SmartDisk{}

	sModel := reSmartModel.FindStringSubmatch(data)
	if len(sModel) > 2 {
		sd.Model = strings.TrimSpace(sModel[1])
		device := strings.TrimSpace(sModel[2])
		if device != "" {
			sd.Model += " (" + device + ")"
		}
	} else {
		sModel = reSmartModel2.FindStringSubmatch(data)
		if len(sModel) > 1 {
			sd.Model = strings.TrimSpace(sModel[1])
		} else {
			sModel = reSmartModel3.FindStringSubmatch(data)
			if len(sModel) > 1 {
				sd.Model = strings.TrimSpace(sModel[1])
			}
		}
	}

	sSerial := reSmartSerial.FindStringSubmatch(data)
	if len(sSerial) > 1 {
		sd.Serial = strings.TrimSpace(sSerial[1])
	}

	sCapacity := reSmartCapacity.FindStringSubmatch(data)
	if len(sCapacity) > 1 {
		sd.Capacity = strings.TrimSpace(sCapacity[1])
	} else {
		sCapacity = reSmartCapacity2.FindStringSubmatch(data)
		if len(sCapacity) > 1 {
			sd.Capacity = strings.TrimSpace(sCapacity[1])
		} else {
			sCapacity = reSmartCapacity3.FindStringSubmatch(data)
			if len(sCapacity) > 1 {
				sd.Capacity = strings.TrimSpace(sCapacity[1])
			}
		}
	}

	sHealth := reSmartHealth.FindStringSubmatch(data)
	if len(sHealth) > 1 {
		sd.Health = strings.TrimSpace(strings.ToUpper(sHealth[1]))
		if sHealth[1] == "PASSED" {
			sd.Health = "OK"
		}
	} else {
		sHealth = reSmartHealth2.FindStringSubmatch(data)
		if len(sHealth) > 1 {
			sd.Health = strings.TrimSpace(strings.ToUpper(sHealth[1]))
			if sHealth[1] == "PASSED" {
				sd.Health = "OK"
			}
		}
	}

	sUsed := reSmartUsed.FindStringSubmatch(data)
	if len(sUsed) > 1 {
		// ssd nvme
		sd.Used = strings.TrimSpace(sUsed[1])
	} else {
		// try to found attr 231 for intel ssd and other
		if attrMatch := reSmartUsedAttr231.FindStringSubmatch(data); len(attrMatch) > 1 {
			if val, err := strconv.Atoi(attrMatch[1]); err == nil && val <= 100 {
				sd.Used = fmt.Sprintf("%d%%", 100-val)
			}
		} else if attrMatch = reSmartUsedAttr233.FindStringSubmatch(data); len(attrMatch) > 1 {
			if val, err := strconv.Atoi(attrMatch[1]); err == nil && val <= 100 {
				sd.Used = fmt.Sprintf("%d%%", 100-val)
			}
		}
	}

	sWorking := reSmartWorking.FindStringSubmatch(data)
	if len(sWorking) > 1 {
		sWorking[1] = strings.Replace(sWorking[1], ",", "", -1)
		if val, err := strconv.Atoi(strings.TrimSpace(sWorking[1])); err == nil {
			sd.Working = int64(val) * 3600
		}
	} else {
		sWorking = reSmartWorking2.FindStringSubmatch(data)
		if len(sWorking) > 1 {
			sWorking[1] = strings.Replace(sWorking[1], ",", "", -1)
			if val, err := strconv.Atoi(strings.TrimSpace(sWorking[1])); err == nil {
				sd.Working = int64(val) * 3600
			}
		}
	}

	sTemperature := reSmartTemperature.FindStringSubmatch(data)
	if len(sTemperature) > 1 {
		sTemperature[1] = strings.Replace(sTemperature[1], "Min/Max ", "", 1)
		sd.Temperature = strings.TrimSpace(sTemperature[1])
	} else {
		sTemperature = reSmartTemperature2.FindStringSubmatch(data)
		if len(sTemperature) > 1 {
			sTemperature[1] = strings.Replace(sTemperature[1], "Min/Max ", "", 1)
			sd.Temperature = strings.TrimSpace(sTemperature[1])
		} else {
			for k, sTemperature := range reSmartTemperature3.FindAllStringSubmatch(data, -1) {
				if len(sTemperature) > 1 {
					sep := ""
					if sd.Temperature != "" {
						sep = " | "
					}
					sd.Temperature += sep + fmt.Sprintf("#%d: ", k) + strings.TrimSpace(sTemperature[1])
				}
			}

			if sd.Temperature == "" {
				sTemperature = reSmartTemperature4.FindStringSubmatch(data)
				if len(sTemperature) > 1 {
					sd.Temperature = strings.TrimSpace(sTemperature[1])
				}
			}
		}
	}

	return sd
}

var reProcMdStat = regexp.MustCompile(`(?s)(md.*?) :(.*?)\n\n`)
var reProcMdDisks = regexp.MustCompile(`(?s)(\S*)\[\d*]`)
var reProcMdState = regexp.MustCompile(`(?s)(check|resync|recovery) = .*?(.+)%.*?finish=(.*)min speed=(.*)K/sec`)
var reMdAdm = regexp.MustCompile(`(?s)Creation Time : (.*?)\n.*?Raid Level : (.*?)\n.*?Array Size : .*?\(.*?B (.*?)\)\n.*?State :(.*?)\n.*?Active Devices : (.*?)\n.*?Working Devices : (.*?)\n.*?Failed Devices : (.*?)\n.*?Spare Devices : (.*?)\n.*?Name : (.*?)(?:\n|\s)`)

func (c *Collector) parseMdAdm(data string) RaidMDAdm {
	adm := RaidMDAdm{}
	m := reMdAdm.FindStringSubmatch(data)
	if len(m) < 7 {
		log.Println("warn: reMdAdm len incorrect")
		return adm
	}
	adm = RaidMDAdm{
		Name:      strings.TrimSpace(m[9]),
		State:     strings.TrimSpace(strings.ToLower(m[4])),
		Level:     strings.TrimSpace(strings.ToUpper(m[2])),
		Capacity:  strings.TrimSpace(m[3]),
		CreatedAt: strings.TrimSpace(m[1]),
		Active:    strings.TrimSpace(m[5]),
		Working:   strings.TrimSpace(m[6]),
		Failed:    strings.TrimSpace(m[7]),
		Spare:     strings.TrimSpace(m[8]),
	}
	return adm
}

func (c *Collector) parseMdStat(data string) map[string]*RaidMD {
	res := map[string]*RaidMD{}
	for _, proc := range reProcMdStat.FindAllStringSubmatch(data, -1) {
		if len(proc) < 3 {
			log.Println("warn: reProcMdStat len incorrect")
			continue
		}
		var disks []string
		for _, md := range reProcMdDisks.FindAllStringSubmatch(proc[2], -1) {
			if len(md) < 2 {
				log.Println("warn: reProcMdDisks len incorrect")
				continue
			}
			disks = append(disks, md[1])
		}

		state := RaidProc{}
		st := reProcMdState.FindStringSubmatch(proc[2])
		if len(st) > 4 {
			state.State = strings.TrimSpace(st[1])
			if state.State != "check" && state.State != "resync" && state.State != "recovery" {
				state.ParseErr = "state: not check, resync or recovery"
				continue
			}
			if progress, err := strconv.ParseFloat(strings.TrimSpace(st[2]), 64); err == nil {
				state.Progress = progress
			} else {
				state.ParseErr = "progress: " + err.Error()
			}
			if left, err := strconv.ParseFloat(strings.TrimSpace(st[3]), 64); err == nil {
				state.Left = left * 60
			} else {
				state.ParseErr = "left: " + err.Error()
			}
			if speed, err := strconv.Atoi(strings.TrimSpace(st[4])); err == nil {
				state.Speed = speed
			} else {
				state.ParseErr = "speed: " + err.Error()
			}
		}

		if logger.IsDebugMode() {
			logger.Debugf("[collector] md proc raw: %s", proc[2])
		}

		res[proc[1]] = &RaidMD{
			Disks: disks,
			Proc:  state,
		}
	}
	return res
}

type ZfsOut struct {
	OutputVersion struct {
		Command   string `json:"command"`
		VersMajor int    `json:"vers_major"`
		VersMinor int    `json:"vers_minor"`
	} `json:"output_version"`
	Pools map[string]json.RawMessage `json:"pools"`
}

type ZfsListPool struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	State      string `json:"state"`
	PoolGuid   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	Properties struct {
		Size struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"size"`
		Allocated struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"allocated"`
		Free struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"free"`
		Checkpoint struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"checkpoint"`
		Expandsize struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"expandsize"`
		Fragmentation struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"fragmentation"`
		Capacity struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"capacity"`
		Dedupratio struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"dedupratio"`
		Health struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"health"`
		Altroot struct {
			Value  string `json:"value"`
			Source struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"source"`
		} `json:"altroot"`
	} `json:"properties"`

	Vdevs map[string]struct {
		Name       string `json:"name"`
		VdevType   string `json:"vdev_type"`
		Guid       string `json:"guid"`
		Class      string `json:"class"`
		State      string `json:"state"`
		Properties struct {
			Size struct {
				Value  string `json:"value"`
				Source struct {
					Type string `json:"type"`
					Data string `json:"data"`
				} `json:"source"`
			} `json:"size"`
			Allocated struct {
				Value  string `json:"value"`
				Source struct {
					Type string `json:"type"`
					Data string `json:"data"`
				} `json:"source"`
			} `json:"allocated"`
			Free struct {
				Value  string `json:"value"`
				Source struct {
					Type string `json:"type"`
					Data string `json:"data"`
				} `json:"source"`
			} `json:"free"`
			Health struct {
				Value  string `json:"value"`
				Source struct {
					Type string `json:"type"`
					Data string `json:"data"`
				} `json:"source"`
			} `json:"health"`
		} `json:"properties"`
		Vdevs map[string]struct {
			Name       string `json:"name"`
			VdevType   string `json:"vdev_type"`
			Guid       string `json:"guid"`
			Path       string `json:"path"`
			Class      string `json:"class"`
			State      string `json:"state"`
			Properties struct {
				Health struct {
					Value  string `json:"value"`
					Source struct {
						Type string `json:"type"`
						Data string `json:"data"`
					} `json:"source"`
				} `json:"health"`
			} `json:"properties"`
		} `json:"vdevs"`
	} `json:"vdevs"`
}

type ZfsStatusPool struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	PoolGuid   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	Status     string `json:"status"`
	Action     string `json:"action"`
	ScanStats  struct {
		Function           string `json:"function"`
		State              string `json:"state"`
		StartTime          string `json:"start_time"`
		EndTime            string `json:"end_time"`
		ToExamine          string `json:"to_examine"`
		Examined           string `json:"examined"`
		Skipped            string `json:"skipped"`
		Processed          string `json:"processed"`
		Errors             string `json:"errors"`
		BytesPerScan       string `json:"bytes_per_scan"`
		PassStart          string `json:"pass_start"`
		ScrubPause         string `json:"scrub_pause"`
		ScrubSpentPaused   string `json:"scrub_spent_paused"`
		IssuedBytesPerScan string `json:"issued_bytes_per_scan"`
		Issued             string `json:"issued"`
	} `json:"scan_stats"`
	Vdevs map[string]struct {
		Name           string `json:"name"`
		VdevType       string `json:"vdev_type"`
		Guid           string `json:"guid"`
		Class          string `json:"class"`
		State          string `json:"state"`
		AllocSpace     string `json:"alloc_space"`
		TotalSpace     string `json:"total_space"`
		DefSpace       string `json:"def_space"`
		ReadErrors     string `json:"read_errors"`
		WriteErrors    string `json:"write_errors"`
		ChecksumErrors string `json:"checksum_errors"`
		Vdevs          map[string]struct {
			Name           string          `json:"name"`
			VdevType       string          `json:"vdev_type"`
			Guid           string          `json:"guid"`
			Class          string          `json:"class"`
			State          string          `json:"state"`
			AllocSpace     string          `json:"alloc_space"`
			TotalSpace     string          `json:"total_space"`
			DefSpace       string          `json:"def_space"`
			RepDevSize     string          `json:"rep_dev_size"`
			ReadErrors     string          `json:"read_errors"`
			WriteErrors    string          `json:"write_errors"`
			ChecksumErrors string          `json:"checksum_errors"`
			Vdevs          json.RawMessage `json:"vdevs"`
		}
	} `json:"vdevs"`
	ErrorCount string `json:"error_count"`
}

func (c *Collector) parseZfs(jList, jStatus []byte) ([]RaidZFS, error) {
	if logger.IsDebugMode() {
		logger.Debugf("[collector] zfs list raw: %s", jList)
		logger.Debugf("[collector] zfs status raw: %s", jStatus)
	}

	rds := make([]RaidZFS, 0)

	var zlo ZfsOut
	err := json.Unmarshal(jList, &zlo)
	if err != nil {
		return nil, fmt.Errorf("zfs list json unmarshal err: %v", err)
	}

	var zso ZfsOut
	err = json.Unmarshal(jStatus, &zso)
	if err != nil {
		return nil, fmt.Errorf("zfs status json unmarshal err: %v", err)
	}

	for _, listRaw := range zlo.Pools {
		var zp ZfsListPool
		err = json.Unmarshal(listRaw, &zp)
		if err != nil {
			return nil, fmt.Errorf("zfs list pool json unmarshal err: %v", err)
		}

		var zs *ZfsStatusPool
		statusRaw, okZsPl := zso.Pools[zp.Name]
		if okZsPl {
			err = json.Unmarshal(statusRaw, &zs)
			if err != nil {
				return nil, fmt.Errorf("zfs status pool json unmarshal err: %v", err)
			}
		}

		for _, dev := range zp.Vdevs {
			devs := make([]RaidZFSDevs, 0)
			for _, d := range dev.Vdevs {
				if d.VdevType != "disk" {
					continue
				}
				devs = append(devs, RaidZFSDevs{
					Name:  strings.Replace(d.Name, "/dev/", "", 1),
					State: strings.TrimSpace(strings.ToLower(d.State)),
				})
			}

			var (
				errRead  = 0
				errWrite = 0
				errChsum = 0
			)
			if zs != nil {
				if zsDevGr, ok := zs.Vdevs[zs.Name]; ok {
					if pDev, ok := zsDevGr.Vdevs[dev.Name]; ok {
						errRead, _ = strconv.Atoi(pDev.ReadErrors)
						errWrite, _ = strconv.Atoi(pDev.WriteErrors)
						errChsum, _ = strconv.Atoi(pDev.ChecksumErrors)
					}
				}
			}

			size, _ := strconv.Atoi(dev.Properties.Size.Value)

			rds = append(rds, RaidZFS{
				PoolName:  zp.Name,
				PoolState: strings.TrimSpace(strings.ToLower(zp.State)),
				Name:      dev.Name,
				Type:      dev.VdevType,
				State:     strings.TrimSpace(strings.ToLower(dev.State)),
				ErrRead:   errRead,
				ErrWrite:  errWrite,
				ErrCkSum:  errChsum,
				Capacity:  size,
				Devs:      devs,
			})
		}
	}

	return rds, nil
}
