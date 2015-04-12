package main

import (
	ui "github.com/gizak/termui"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	ramPath string        = "/proc/meminfo"
	delay   time.Duration = 2
)

type MemInfo struct {
	MemTotal uint64
	MemFree  uint64
	MemUsed  uint64
}

func (m *MemInfo) changeMB() {
	m.MemTotal /= 1024
	m.MemFree /= 1024
	m.MemUsed /= 1024
}

func (m *MemInfo) calPercentage() {
	m.MemFree = (m.MemFree * 100) / m.MemTotal
	m.MemUsed = (m.MemUsed * 100) / m.MemTotal
	m.MemTotal = (m.MemTotal * 100) / m.MemTotal //doesn't make any sense
}

func getRAMStat() (*MemInfo, error) {
	var mem = MemInfo{}

	fData, err := ioutil.ReadFile(ramPath)

	if err != nil {
		return nil, err
	}

	sData := string(fData)

	lines := strings.Split(sData, "\n")

	for inx, line := range lines {
		fields := strings.Fields(line)
		if inx == 0 {
			mem.MemTotal, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if inx == 1 {
			mem.MemFree, _ = strconv.ParseUint(fields[1], 10, 64)
		} else {
			break
		}
	}
	mem.MemUsed = mem.MemTotal - mem.MemFree
	mem.changeMB()
	return &mem, nil
}

func memstat() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	ram_ls := ui.NewList()
	ram_ls.HasBorder = false
	ram_ls.Items = []string{
		"TOTAL RAM MEMORY",
		"",
		"FREE RAM MEMORY",
		"",
		"USED RAM MEMORY",
	}
	ram_ls.Height = 5
	ram_ls.Width = 25
	ram_ls.X = 1
	ram_ls.Y = 1

	ram_gs := make([]*ui.Gauge, 3)
	for i := range ram_gs {
		ram_gs[i] = ui.NewGauge()
		ram_gs[i].Height = 2
		ram_gs[i].HasBorder = false
		ram_gs[i].Percent = i * 10
		ram_gs[i].PaddingBottom = 1
		ram_gs[i].BarColor = ui.ColorBlue
		ram_gs[i].Width = 50
		ram_gs[i].X = 25
		ram_gs[i].Y = 1 + i*2
	}

	ramsize_ls := ui.NewList()
	ramsize_ls.HasBorder = false
	ramsize_ls.Items = []string{
		"0MB",
		"",
		"0MB",
		"",
		"0MB",
	}
	ramsize_ls.Height = 5
	ramsize_ls.Width = 25
	ramsize_ls.X = 85
	ramsize_ls.Y = 1

	draw := func() {
		ui.Render(ram_ls, ram_gs[0], ram_gs[1], ram_gs[2], ramsize_ls)
	}

	evt := ui.EventCh()

	var mem *MemInfo

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		default:
			mem, _ = getRAMStat()
			ramsize_ls.Items[0] = strconv.Itoa((int)(mem.MemTotal)) + " MB"
			ramsize_ls.Items[2] = strconv.Itoa((int)(mem.MemFree)) + " MB"
			ramsize_ls.Items[4] = strconv.Itoa((int)(mem.MemUsed)) + " MB"

			mem.calPercentage()

			ram_gs[0].Percent = (int)(mem.MemTotal)
			ram_gs[1].Percent = (int)(mem.MemFree)
			ram_gs[2].Percent = (int)(mem.MemUsed)

			draw()
			time.Sleep(delay * time.Second)
		}
	}

}
