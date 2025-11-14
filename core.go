package main

import (
	"github.com/gorilla/websocket"
	"log"
	tests "netip-core/benchmark"
	"netip-core/collector"
	"netip-core/info"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// reader from nodes-handler
	go func() {
		for {
			if !conn.Alive() {
				time.Sleep(time.Second)
				continue
			}

			var res struct {
				Command string `json:"command"`
				Runtime int    `json:"runtime"`
			}

			// needs for ping-pong
			err := conn.ws.ReadJSON(&res)
			if err != nil {
				if conn.Alive() && !websocket.IsCloseError(err, 1006) {
					log.Println("[component] ws read err:", err)
				}
				time.Sleep(time.Second)
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
			conn.Write(struct {
				Event       string                 `json:"event"`
				CollectCore *collector.CollectCore `json:"collectCore"`
			}{
				Event:       "collect-core",
				CollectCore: cc,
			}, "chan core")

		// chan-sender who logged terminals
		case wl, ok := <-col.ChanWhoLogged:
			if !ok {
				continue
			}
			conn.Write(struct {
				Event     string               `json:"event"`
				WhoLogged *collector.WhoLogged `json:"whoLogged"`
			}{
				Event:     "who-logged",
				WhoLogged: wl,
			}, "who logged")

		// chan-sender processes
		case ps, ok := <-col.ChanProcesses:
			if !ok {
				continue
			}
			conn.Write(struct {
				Event     string               `json:"event"`
				Processes *collector.Processes `json:"processes"`
			}{
				Event:     "processes",
				Processes: ps,
			}, "processes")

		// chan-sender general-test
		case gt, ok := <-chGeneralTests:
			if !ok {
				continue
			}
			conn.Write(struct {
				Event    string        `json:"event"`
				BmsTests *tests.Result `json:"bmsTests"`
			}{
				Event:    "bms-general-tests",
				BmsTests: gt,
			}, "general-tests")

		// chan-sender stats disks info
		case cdi, ok := <-col.ChanDisksInfo:
			if !ok {
				continue
			}
			conn.Write(struct {
				Event     string               `json:"event"`
				DisksInfo *collector.DisksInfo `json:"disksInfo"`
			}{
				Event:     "disks-info",
				DisksInfo: cdi,
			}, "disks-info")

		// handler destroy
		case <-destroy:
			log.Println("[component] service destroyed")
			log.Println("------")
			log.Println("below remains to execution manually:")
			log.Println("# docker rm -f netip.core")
			log.Println("------")
			conn.SendClose()
			<-destroy

		// handler terminate
		case <-terminate:
			log.Println("[component] terminating...")
			conn.SendClose()
			os.Exit(0)
		}
	}
}
