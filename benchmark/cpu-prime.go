package tests

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var testCpuPrime = &TestSoftVer{
	Test:     "cpu-prime",
	Software: "sysbench",
	Version:  "1.0.20",
}

var reBMCpuPrime = regexp.MustCompile(`events per second:\s*([0-9.]+)`)

func BMCpuPrime(ctx context.Context, runtime int) (float64, error) {
	res, err := shell(ctx,
		fmt.Sprintf("sysbench cpu --cpu-max-prime=10000 --time=%d --threads=1 run", runtime))
	if err != nil {
		return 0, err
	}
	match := reBMCpuPrime.FindStringSubmatch(res)
	if len(match) > 1 {
		return strconv.ParseFloat(match[1], 64)
	}
	return 0, errors.New("regex err")
}
