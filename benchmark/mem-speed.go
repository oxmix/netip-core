package tests

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var testMemSpeed = &TestSoftVer{
	Test:     "mem-speed",
	Software: "sysbench",
	Version:  "1.0.20",
}

var reBMMemSpeedMiB = regexp.MustCompile(`transferred \(([0-9.]+) MiB/sec\)`)

func BMMemSpeedMiB(ctx context.Context, runtime int) (float64, error) {
	res, err := shell(ctx,
		fmt.Sprintf("sysbench memory --memory-total-size=8G --threads=1 --time=%d run", runtime))
	if err != nil {
		return 0, err
	}
	match := reBMMemSpeedMiB.FindStringSubmatch(res)
	if len(match) > 1 {
		return strconv.ParseFloat(match[1], 64)
	}
	return 0, errors.New("regex err")
}
