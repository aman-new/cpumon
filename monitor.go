package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Monitor struct {
	fr          FileReader
	cr          CmdRunner
	cpuModel    string
	deviceModel string
	cpuFreqs    []CPUFreqInfo
	hwmonPath   string
	hwmonTemps  []string
	fanFiles    []string
	thinkpadFan bool
	sensorsOK   bool

	coreFreqBuf map[int]string
	lineBuf     []string
}

func NewMonitor() *Monitor {
	fr := sysFileReader{}
	cr := sysCmdRunner{}

	_, sensorsErr := exec.LookPath("sensors")
	hwmonPath := discoverHwmonCPU(fr)

	return &Monitor{
		fr:          fr,
		cr:          cr,
		cpuModel:    readCPUModel(fr),
		deviceModel: readDeviceModel(fr),
		cpuFreqs:    discoverCPUTopology(fr),
		hwmonPath:   hwmonPath,
		hwmonTemps:  discoverHwmonTemps(hwmonPath),
		fanFiles:    discoverFanFiles(),
		thinkpadFan: fileExists(thinkpadFanPath),
		sensorsOK:   sensorsErr == nil,
		coreFreqBuf: make(map[int]string, 32),
		lineBuf:     make([]string, 0, 32),
	}
}

func (m *Monitor) collect() Metrics {
	avgFreq := readFrequencies(m.fr, m.cpuFreqs, m.coreFreqBuf)

	cpuStatus, _ := readCPUThermal(m.fr, m.cr, m.sensorsOK, m.hwmonTemps, m.coreFreqBuf, &m.lineBuf)
	fanStatus, _ := readFanStatus(m.fr, m.fanFiles, m.thinkpadFan, &m.lineBuf)

	return Metrics{
		DeviceModel: m.deviceModel,
		CPUModel:    m.cpuModel,
		Governor:    readGovernor(m.fr),
		EnergyBias:  readEnergyBias(m.fr),
		AvgFreq:     avgFreq,
		CPUStatus:   cpuStatus,
		Throttle:    readThrottleInfo(m.fr),
		FanStatus:   fanStatus,
		SensorsHint: !m.sensorsOK,
	}
}

func (m *Monitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	var wg sync.WaitGroup
	wg.Go(func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	})

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	display(m.collect())

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nShutting down...")
			wg.Wait()
			return nil
		case <-ticker.C:
			display(m.collect())
		}
	}
}
