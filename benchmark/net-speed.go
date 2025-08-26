package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var testNetSpeed = &TestSoftVer{
	Test:     "net-speed",
	Software: "speedtest-cli",
	Version:  "2.1.3",
}

type SpeedTestCli struct {
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
	Ping     float64 `json:"ping"`
	Server   struct {
		Url     string  `json:"url"`
		Lat     string  `json:"lat"`
		Lon     string  `json:"lon"`
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Cc      string  `json:"cc"`
		Sponsor string  `json:"sponsor"`
		Id      string  `json:"id"`
		Host    string  `json:"host"`
		D       float64 `json:"d"`
		Latency float64 `json:"latency"`
	} `json:"server"`
	Timestamp     time.Time `json:"timestamp"`
	BytesSent     int       `json:"bytes_sent"`
	BytesReceived int       `json:"bytes_received"`
	Share         any       `json:"share"`
	Client        struct {
		Ip        string `json:"ip"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Isp       string `json:"isp"`
		Isprating string `json:"isprating"`
		Rating    string `json:"rating"`
		Ispdlavg  string `json:"ispdlavg"`
		Ispulavg  string `json:"ispulavg"`
		Loggedin  string `json:"loggedin"`
		Country   string `json:"country"`
	} `json:"client"`
}

func BMNetSpeed(ctx context.Context, runtime int) (*SpeedTestCli, error) {
	res, err := shell(ctx, fmt.Sprintf("speedtest-cli --timeout %d --json", runtime))
	if err != nil {
		return nil, err
	}
	stc := new(SpeedTestCli)
	err = json.Unmarshal([]byte(res), stc)
	if err != nil {
		return nil, err
	}
	return stc, nil
}
