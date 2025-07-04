package tests

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var lock bool

type TestSoftVer struct {
	Test     string `json:"test"`
	Software string `json:"software"`
	Version  string `json:"version"`
}

type Result struct {
	Test      string `json:"test"`
	Scheduled bool   `json:"scheduled"`
	Runtime   int    `json:"runtime"`
	CPU       struct {
		*TestSoftVer
		EventsSec float64 `json:"eventsSec"`
	} `json:"cpu"`
	Mem struct {
		*TestSoftVer
		SpeedMiB float64 `json:"speedMib"`
	} `json:"mem"`
	IO struct {
		*TestSoftVer
		*FioResult
	} `json:"io"`
	Net struct {
		*TestSoftVer
		*SpeedTestCli
	} `json:"net"`
}

func NewGeneralTests(scheduled bool, runtime int, channel chan<- *Result) {
	if lock {
		log.Println("locked by other test")
		return
	}
	lock = true
	defer func() {
		lock = false
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	out := &Result{
		Test:      "general-tests",
		Scheduled: scheduled,
		Runtime:   runtime,
	}

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		events, err := BMCpuPrime(ctx, runtime)
		if err != nil {
			log.Println("general-tests cpu err:", err)
			return
		}
		out.CPU.TestSoftVer = testCpuPrime
		out.CPU.EventsSec = events
	}()

	go func() {
		defer wg.Done()
		speedMiB, err := BMMemSpeedMiB(ctx, runtime)
		if err != nil {
			log.Println("general-tests mem err:", err)
			return
		}
		out.Mem.TestSoftVer = testMemSpeed
		out.Mem.SpeedMiB = speedMiB
	}()

	go func() {
		defer wg.Done()
		fioRes, err := BMFio(ctx, runtime)
		if err != nil {
			log.Println("general-tests io err:", err)
			return
		}
		out.IO.TestSoftVer = testFio
		out.IO.FioResult = fioRes
	}()

	go func() {
		defer wg.Done()
		netRes, err := BMNetSpeed(ctx, runtime)
		if err != nil {
			log.Println("general-tests net err:", err)
			return
		}
		out.Net.TestSoftVer = testNetSpeed
		out.Net.SpeedTestCli = netRes
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			log.Println("general-tests: interrupted due to timeout")
			return

		case <-done:
			log.Println("general-tests: ok")
			channel <- out
			return
		}
	}
}

func shell(ctx context.Context, command string) (string, error) {
	out, err := exec.CommandContext(ctx, "sh", "-c", command).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("shell err: %s command: %s", err, command)
	}
	return strings.TrimSpace(string(out)), nil
}
