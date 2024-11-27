package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var testFio = &TestSoftVer{
	Test:     "io",
	Software: "fio",
	Version:  "3.34",
}

type FioPrepare struct {
	Jobs []struct {
		Jobname string `json:"jobname"`
		Error   int    `json:"error"`
		Read    struct {
			BwBytes uint    `json:"bw_bytes"` // KiB/s=1024 684,795,686
			Iops    float64 `json:"iops"`     // "iops" : 167186.446893,
			LatNs   struct {
				Mean float64 `json:"mean"` // (latency) в наносекундах (ns)
			} `json:"lat_ns"`
		} `json:"read"`
		Write struct {
			BwBytes uint    `json:"bw_bytes"`
			Iops    float64 `json:"iops"`
			LatNs   struct {
				Mean float64 `json:"mean"` // (latency) в наносекундах (ns)
			} `json:"lat_ns"`
		} `json:"write"`
		Trim struct {
			BwBytes uint    `json:"bw_bytes"`
			Iops    float64 `json:"iops"`
			LatNs   struct {
				Mean float64 `json:"mean"`
			} `json:"lat_ns"`
		} `json:"trim"`
		Sync struct {
			TotalIos uint `json:"total_ios"`
			LatNs    struct {
				Mean float64 `json:"mean"`
			} `json:"lat_ns"`
		} `json:"sync"`
	} `json:"jobs"`
}

type FioResult struct {
	ReadBw     uint    `json:"read"`
	ReadIops   float64 `json:"readIops"`
	ReadLatNs  float64 `json:"readLatNs"`
	WriteBw    uint    `json:"write"`
	WriteIops  float64 `json:"writeIops"`
	WriteLatNs float64 `json:"writeLatNs"`
}

func BMFio(ctx context.Context, runtime int) (*FioResult, error) {
	res, err := shell(ctx, "fio --output-format=json --name=io_test --rw=readwrite --bs=4k"+
		fmt.Sprintf(" --size=1G --numjobs=1 --runtime=%d --time_based --group_reporting", runtime))
	if err != nil {
		return nil, err
	}
	var fio FioPrepare
	err = json.Unmarshal([]byte(res), &fio)
	if err != nil {
		return nil, err
	}
	if len(fio.Jobs) == 0 {
		return nil, errors.New("fio output is empty")
	}
	if fio.Jobs[0].Error > 0 {
		return nil, fmt.Errorf("fio jobs err: %d", fio.Jobs[0].Error)
	}

	return &FioResult{
		ReadBw:     fio.Jobs[0].Read.BwBytes,
		ReadIops:   fio.Jobs[0].Read.Iops,
		ReadLatNs:  fio.Jobs[0].Read.LatNs.Mean,
		WriteBw:    fio.Jobs[0].Write.BwBytes,
		WriteIops:  fio.Jobs[0].Write.Iops,
		WriteLatNs: fio.Jobs[0].Write.LatNs.Mean,
	}, nil
}
