package main

import (
	"fmt"
	"strconv"
)

type ThrottleInfo struct {
	Available        bool
	PackageCount     string
	PackageTotalTime string
	PackageMaxTime   string
	CoreCount        string
	CoreTotalTime    string
	CoreMaxTime      string
}

func readThrottleInfo(fr FileReader) ThrottleInfo {
	if !fileExists(cpuThrottlePath) {
		return ThrottleInfo{Available: false}
	}

	info := ThrottleInfo{Available: true}

	info.PackageCount = readCount(fr, cpuThrottlePath)
	info.PackageTotalTime = readDuration(fr, pkgThrottleTotal)
	info.PackageMaxTime = readMs(fr, pkgThrottleMax)
	info.CoreCount = readCount(fr, coreThrottleCount)
	info.CoreTotalTime = readDuration(fr, coreThrottleTotal)
	info.CoreMaxTime = readMs(fr, coreThrottleMax)

	return info
}

func readCount(fr FileReader, path string) string {
	raw, err := fr.Read(path)
	if err != nil || raw == "" {
		return "N/A"
	}
	if _, err := strconv.ParseInt(raw, 10, 64); err != nil {
		return "N/A"
	}
	return raw
}

func readMs(fr FileReader, path string) string {
	raw, err := fr.Read(path)
	if err != nil || raw == "" {
		return "N/A"
	}
	ms, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%d ms", ms)
}

func readDuration(fr FileReader, path string) string {
	raw, err := fr.Read(path)
	if err != nil || raw == "" {
		return "N/A"
	}
	ms, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return "N/A"
	}
	return formatDuration(ms)
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	sec := ms / 1000
	if sec < 60 {
		return fmt.Sprintf("%d seconds", sec)
	}
	min := sec / 60
	remSec := sec % 60
	if min < 60 {
		return fmt.Sprintf("%d min %d sec", min, remSec)
	}
	hr := min / 60
	remMin := min % 60
	return fmt.Sprintf("%d hr %d min", hr, remMin)
}
