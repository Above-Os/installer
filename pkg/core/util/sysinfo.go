package util

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetHost() {
	hostInfo, _ := host.Info()
	fmt.Printf("host info hostName: %s, hostId: %s, os: %s, platform: %s, version: %s, arch: %s\n",
		hostInfo.Hostname, hostInfo.HostID, hostInfo.OS, hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch)
}

func GetCpu() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cpuInfo, _ := cpu.InfoWithContext(ctx)
	if len(cpuInfo) == 0 {
		return
	}

	cpuLogicalCount, _ := cpu.CountsWithContext(ctx, true)
	cpuPhysicalCount, _ := cpu.CountsWithContext(ctx, false)

	fmt.Printf("cpu info model: %s, logicalCount: %d, physicalCount: %d\n", cpuInfo[0].ModelName, cpuLogicalCount, cpuPhysicalCount)
}

func GetDisk() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	usageInfo, _ := disk.UsageWithContext(ctx, "/")
	fmt.Printf("disk info usage total: %d, free: %d\n", usageInfo.Total, usageInfo.Free)
}

func GetMem() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	memInfo, _ := mem.SwapMemoryWithContext(ctx)
	fmt.Printf("mem info total: %d, free: %d, percent: %f\n", memInfo.Total, memInfo.Free, memInfo.UsedPercent)

}
