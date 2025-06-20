package info

import (
	"os"
	"runtime"
	"strings"
)

func (i *Info) fillKernelInfo() {
	i.Data.Kernel.Architecture = runtime.GOARCH
	i.Data.Kernel.OSType = runtime.GOOS

	osRelease, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err == nil {
		i.Data.Kernel.OSRelease = strings.TrimSpace(string(osRelease))
	}

	osVersion, err := os.ReadFile("/proc/sys/kernel/version")
	if err == nil {
		i.Data.Kernel.OSVersion = strings.TrimSpace(string(osVersion))
	}
}
