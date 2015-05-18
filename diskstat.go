package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type statIO struct {
	Name        string
	ReadSector  uint64
	WriteSector uint64
	ReadBytes   uint64
	WriteBytes  uint64
}

const (
	procPath string = "/proc/diskstats"
)

func getDiskstat() (stat []statIO) {
	stat = []statIO{}

	fData, err := ioutil.ReadFile(procPath)

	if err != nil {
		fmt.Println("failed to open", procPath)
		return nil
	}

	strData := string(fData)

	lines := strings.Split(strData, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if fields[3] == "0" {
			continue
		}

		temp_stat := statIO{}
		temp_stat.Name = fields[2]
		temp_stat.ReadSector, _ = strconv.ParseUint(fields[5], 10, 64)
		temp_stat.WriteSector, _ = strconv.ParseUint(fields[9], 10, 64)
		temp_stat.ReadBytes = temp_stat.ReadSector * 512   //Sector to bytes
		temp_stat.WriteBytes = temp_stat.WriteSector * 512 //Sector to bytes

		stat = append(stat, temp_stat)
	}

	return stat
}

func diskstat() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	par0 := ui.NewPar("Disk IO Stats")
	par0.Height = 1
	par0.Width = 50
	par0.X = 20
	par0.Y = 1
	par0.HasBorder = false

	par1 := ui.NewPar("          0 kb/s ")
	par1.Height = 3
	par1.Width = 50
	par1.Y = 3
	par1.Border.Label = "Total Reads"
	par1.Border.FgColor = ui.ColorYellow

	par2 := ui.NewPar("          0 mb/s")
	par2.Height = 3
	par2.Width = 50
	par2.Y = 6
	par2.Border.Label = "Total Writes"
	par2.Border.FgColor = ui.ColorYellow

	draw := func() {
		ui.Render(par0, par1, par2)
	}

	evt := ui.EventCh()

	var prev_read, prev_write, cur_read, cur_write uint64

	stat := getDiskstat()
	for _, v := range stat {
		prev_read += v.ReadBytes
		prev_write += v.WriteBytes
	}

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		default:
			time.Sleep(1 * time.Second)

			stat = getDiskstat()
			for _, v := range stat {
				cur_read += v.ReadBytes
				cur_write += v.WriteBytes
			}
			read_text := fmt.Sprintf("     %d kb / s", (cur_read-prev_read)/1024)
			write_text := fmt.Sprintf("     %d kb / s", (cur_write-prev_write)/1024)
			//read_text := strconv.Itoa(int((cur_read-prev_read)/1024)) + " Kb / s"
			//write_text := strconv.Itoa(int((cur_write-prev_write)/1024)) + " Kb / s"
			prev_read = cur_read
			prev_write = cur_write
			cur_read, cur_write = 0, 0

			par1.Text = read_text
			par2.Text = write_text

			draw()
		}

	}
}
