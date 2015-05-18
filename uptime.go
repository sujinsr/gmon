package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	proc_path string = "/proc/uptime"
)

func getUptime() float64 {
	var utime float64

	file_data, err := ioutil.ReadFile(proc_path)
	if err != nil {
		panic(err)
	}
	str_data := string(file_data)

	lines := strings.Split(str_data, "\n")
	fields := strings.Fields(lines[0])

	utime, _ = strconv.ParseFloat(fields[0], 64)

	return utime
}

func uptime() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	par0 := ui.NewPar("System Up Time")
	par0.Height = 1
	par0.Width = 20
	par0.Y = 2
	par0.X = 20
	par0.HasBorder = false

	par1 := ui.NewPar("")
	par1.Height = 3
	par1.Width = 55
	par1.X = 5
	par1.Y = 3
	par1.Border.FgColor = ui.ColorYellow

	t := getUptime()
	now := time.Now()
	up_sec := now.Unix() - int64(t)

	if t < 60 {
		par1.Text = fmt.Sprintf("%.2f secs Before At %s", getUptime(), time.Unix(up_sec, 0))
	} else if t >= 60 && t < 60*60 {
		par1.Text = fmt.Sprintf("%.2f mins Before At %s", getUptime()/60, time.Unix(up_sec, 0))
	} else {
		par1.Text = fmt.Sprintf("%.2f hrs Before At %s", getUptime()/(60*60), time.Unix(up_sec, 0))
	}

	ui.Render(par0, par1)

	<-ui.EventCh()

}
