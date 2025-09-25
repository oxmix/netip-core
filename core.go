package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	tests "netip-core/benchmark"
	"netip-core/collector"
	"netip-core/info"
	"os"
	"os/signal"
	"time"
)

type connectResponse struct {
	responseBase
}

type connectPayload struct {
	payloadBase
	Info *info.Info `json:"info"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	destroy := make(chan bool, 1)

	connect(&connectPayload{
		payloadBase: payloadBase{
			Service: "core",
		},
		Info: info.Get(),
	})

	col := collector.New()

	chGeneralTests := make(chan *tests.Result)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// reader from nodes-handler
	go func() {
		for {
			if wsConnect == nil {
				time.Sleep(time.Second)
				continue
			}

			// ReadMessage needs for ping-pong
			code, message, err := wsConnect.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, 1006) {
					log.Println("ws read err:", err, code)
				}
				time.Sleep(time.Second)
				continue
			}
			var res struct {
				Command string `json:"command"`
				Runtime int    `json:"runtime"`
			}
			err = json.Unmarshal(message, &res)
			if err != nil {
				log.Println("ws rx unmarshal err:", err)
				continue
			}

			switch res.Command {
			case "general-tests":
				go tests.NewGeneralTests(false, res.Runtime, chGeneralTests)
			case "services-destroy":
				destroy <- true
				return
			}
		}
	}()

	log.Println("ready to work")

	// writer to nodes-handler
	for {
		select {
		case <-destroy:
			log.Println("service destroyed")
			log.Println("------")
			log.Println("below remains to execution manually:")
			log.Println("# docker rm -f netip.core")
			log.Println("------")
			_ = wsConnect.Close()
			<-destroy

		// ticker-sender stats core
		case <-ticker.C:
			if wsConnect == nil {
				continue
			}
			j := struct {
				Event       string                `json:"event"`
				CollectCore collector.CollectCore `json:"collectCore"`
			}{
				Event:       "collect-core",
				CollectCore: col.Get(),
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				// ! lock this select before established connection
				connectDegrade(errors.New("ticker write err: " + err.Error()))
				continue
			}

		// chan-sender who logged terminals
		case wl, ok := <-col.ChanWhoLogged:
			if wsConnect == nil || !ok {
				continue
			}
			j := struct {
				Event     string               `json:"event"`
				WhoLogged *collector.WhoLogged `json:"whoLogged"`
			}{
				Event:     "who-logged",
				WhoLogged: wl,
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				log.Println("who logged write err:", err)
				continue
			}

		// chan-sender processes
		case ps, ok := <-col.ChanProcesses:
			if wsConnect == nil || !ok {
				continue
			}
			j := struct {
				Event     string               `json:"event"`
				Processes *collector.Processes `json:"processes"`
			}{
				Event:     "processes",
				Processes: ps,
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				log.Println("processes write err:", err)
				continue
			}

		// chan-sender general-test
		case gt, ok := <-chGeneralTests:
			if wsConnect == nil || !ok {
				continue
			}
			j := struct {
				Event    string        `json:"event"`
				BmsTests *tests.Result `json:"bmsTests"`
			}{
				Event:    "bms-general-tests",
				BmsTests: gt,
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				log.Println("general-test write err:", err)
				continue
			}

		// chan-sender stats disks info
		case cdi, ok := <-col.ChanDisksInfo:
			if wsConnect == nil || !ok {
				continue
			}
			j := &struct {
				Event     string               `json:"event"`
				DisksInfo *collector.DisksInfo `json:"disksInfo"`
			}{
				Event:     "disks-info",
				DisksInfo: cdi,
			}
			err := wsConnect.WriteJSON(j)
			if err != nil {
				log.Println("disks info write err:", err)
				continue
			}

		// handler terminate
		case <-interrupt:
			log.Println("interrupt")
			if wsConnect == nil {
				return
			}

			err := wsConnect.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close err:", err)
				return
			}
			err = wsConnect.Close()
			if err != nil {
				log.Println("close conn err:", err)
			}
			return
		}
	}
}
