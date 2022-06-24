package proc

import (
	"os"
	"strings"
)

// Kernel contains information about the kernel.
type Kernel struct {
	OSType      string
	OSRelease   string
	BootArgs    string
	FullVersion string
}

// GetKernelInfo returns information about the kernel.
func GetKernelInfo() Kernel {

	var kernel Kernel

	fullVersion, _ := os.ReadFile("/proc/version")
	kernel.FullVersion = strings.TrimSpace(string(fullVersion))

	osType, _ := os.ReadFile("/proc/sys/kernel/ostype")
	kernel.OSType = strings.TrimSpace(string(osType))

	osRelease, _ := os.ReadFile("/proc/sys/kernel/osrelease")
	kernel.OSRelease = strings.TrimSpace(string(osRelease))

	bootArgs, _ := os.ReadFile("/proc/cmdline")
	kernel.BootArgs = strings.TrimSpace(string(bootArgs))

	return kernel
}
