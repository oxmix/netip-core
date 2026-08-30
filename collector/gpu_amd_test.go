package collector

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestParseGpuAmd(t *testing.T) {
	t.Parallel()

	// gpu_metrics v2.2
	data, _ := parseHexDump(`
0000000 0080 0202 0fa0 0fa0 0f55 0f55 0f55 0f3c
0000010 0f6e 0f55 0f87 0f3c 0f6e 0000 0000 0000
0000020 5053 f315 4299 0008 0007 13e1 02b5 ffff
0000030 0020 02fe 0003 0000 0000 0004 0003 0009
0000040 0190 0190 ffff 02f3 0190 ffff 0190 0190
0000050 0006 0190 0190 0190 0af0 0af0 08c0 0000
0000060 0000 08c0 08c0 08c0 0af0 0000 0000 0000
0000070 0000 ffff ffff ffff 0000 0000 0000 0000`)

	ga := Collector{}
	amd, err := ga.parseGpuAmd(0, data)
	if err != nil {
		t.Fatal(err)
	}
	if amd.VerFormat != 2 || amd.VerContent != 2 || amd.Power != "6.93" || amd.ClockSoc != 400 {
		t.Fatalf("error parse gpu amd v2.2: %+v", amd)
	}
}

func parseHexDump(dump string) ([]byte, error) {
	var res []byte
	for _, line := range strings.Split(dump, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		for _, w := range fields[1:] {
			b, err := hex.DecodeString(w)
			if err != nil {
				return nil, err
			}
			res = append(res, b[1], b[0])
		}
	}
	return res, nil
}
