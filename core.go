package main

import (
	"encoding/json"
	"log"
	tests "netip-core/benchmark"
	"netip-core/collector"
	"netip-core/info"
	"os"
	"os/signal"
	"syscall"
)

type ConnectResponse struct {
	ResponseBase
}

type ConnectPayload struct {
	PayloadBase
	Info *info.Info `json:"info"`
}

func main() {
	go collector.Pprof()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	destroy := make(chan struct{}, 1)

	conn := NewConnection(&ConnectPayload{
		PayloadBase: PayloadBase{
			Service: "core",
		},
		Info: info.Get(),
	})

	col := collector.New()
	chGeneralTests := make(chan *tests.Result, 1)

	// live from nodes-handler
	go func() {
		for p := range conn.chanLive {
			var res struct {
				Command string `json:"command"`
				Runtime int    `json:"runtime"`
			}
			err := json.Unmarshal(p, &res)
			if err != nil {
				log.Println("[component] decode live err:", err)
				continue
			}

			switch res.Command {
			case "general-tests":
				go tests.NewGeneralTests(false, res.Runtime, chGeneralTests)
			case "services-destroy":
				destroy <- struct{}{}
				return
			}
		}
	}()

	log.Println("[component] ready to work")

	// writer to nodes-handler
	for {
		select {
		// chan-sender stats core
		case cc, ok := <-col.ChanCore:
			if !ok {
				continue
			}
			conn.chanSend <- struct {
				Event       string                 `json:"event"`
				CollectCore *collector.CollectCore `json:"collectCore"`
			}{
				Event:       "collect-core",
				CollectCore: cc,
			}

		// chan-sender who logged terminals
		case wl, ok := <-col.ChanWhoLogged:
			if !ok {
				continue
			}
			conn.chanSend <- struct {
				Event     string               `json:"event"`
				WhoLogged *collector.WhoLogged `json:"whoLogged"`
			}{
				Event:     "who-logged",
				WhoLogged: wl,
			}

		// chan-sender processes
		case ps, ok := <-col.ChanProcesses:
			if !ok {
				continue
			}
			conn.chanSend <- struct {
				Event     string               `json:"event"`
				Processes *collector.Processes `json:"processes"`
			}{
				Event:     "processes",
				Processes: ps,
			}

		// chan-sender general-test
		case gt, ok := <-chGeneralTests:
			if !ok {
				continue
			}
			conn.chanSend <- struct {
				Event    string        `json:"event"`
				BmsTests *tests.Result `json:"bmsTests"`
			}{
				Event:    "bms-general-tests",
				BmsTests: gt,
			}

		// chan-sender stats disks info
		case cdi, ok := <-col.ChanDisksInfo:
			if !ok {
				continue
			}
			conn.chanSend <- struct {
				Event     string               `json:"event"`
				DisksInfo *collector.DisksInfo `json:"disksInfo"`
			}{
				Event:     "disks-info",
				DisksInfo: cdi,
			}

		// handler destroy
		case <-destroy:
			log.Println("[component] service destroyed")
			log.Println("------")
			log.Println("below remains to execution manually:")
			log.Println("# docker rm -f netip.core")
			log.Println("------")
			conn.close()
			<-destroy

		// handler terminate
		case <-terminate:
			log.Println("[component] terminating...")
			conn.close()
			os.Exit(0)
		}
	}
}
