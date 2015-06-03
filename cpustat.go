package main

import (
	ui "github.com/gizak/termui"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type rawStat struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

type StatCPU struct {
	User    float32
	Nice    float32
	System  float32
	Idle    float32
	Iowait  float32
	Irq     float32
	SoftIrq float32
	Steal   float32
	Guest   float32
}

func parseCPUFields(fields []string, stat *rawStat) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}

var preCPU rawStat

func getCPU(stats *StatCPU) (err error) {
	var nowCPU rawStat
	var total float32

	fData, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}

	sData := string(fData)

	lines := strings.Split(sData, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			parseCPUFields(fields, &nowCPU)
			break
		}

	}

	if preCPU.Total != 0 { 
		total = float32(nowCPU.Total - preCPU.Total)
		stats.User = float32(nowCPU.User-preCPU.User) / total * 100
		stats.Nice = float32(nowCPU.Nice-preCPU.Nice) / total * 100
		stats.System = float32(nowCPU.System-preCPU.System) / total * 100
		stats.Idle = float32(nowCPU.Idle-preCPU.Idle) / total * 100
		stats.Iowait = float32(nowCPU.Iowait-preCPU.Iowait) / total * 100
		stats.Irq = float32(nowCPU.Irq-preCPU.Irq) / total * 100
		stats.SoftIrq = float32(nowCPU.SoftIrq-preCPU.SoftIrq) / total * 100
		stats.Guest = float32(nowCPU.Guest-preCPU.Guest) / total * 100
	}
	preCPU = nowCPU
	return 
}

func loadStat(stats *StatCPU, data []int) {
	data[0] = int(stats.User)
	data[1] = int(stats.Nice)
	data[2] = int(stats.System)
	data[3] = int(stats.Idle)
	data[4] = int(stats.Iowait)
	data[5] = int(stats.Irq)
	data[6] = int(stats.SoftIrq)
	data[7] = int(stats.Steal)
	data[8] = int(stats.Guest)

}

func cpustat() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	bc := ui.NewBarChart()
	bclabels := []string{"user", "nice", "sys", "idle", "iow", "irq", "s_irq", "steal", "guest"}
	bc.Border.Label = "CPU Stat"
	bc.X = 15
	bc.Y = 2
	bc.Width = 75
	bc.Height = 30
	bc.BarWidth = 5
	bc.BarGap = 3
	bc.DataLabels = bclabels
	bc.TextColor = ui.ColorGreen
	bc.BarColor = ui.ColorRed
	bc.NumColor = ui.ColorYellow

	stats := StatCPU{}
	data := make([]int, 9)

	evt := ui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		default:
			getCPU(&stats)
			loadStat(&stats, data)
			bc.Data = data
			ui.Render(bc)
			time.Sleep(1 * time.Second)
		}
	}
}
