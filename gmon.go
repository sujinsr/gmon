package main

import (
	"fmt"
	"os"
)

func main() {
	/* get command line arguments without program name */
	Args := os.Args[1:]

	if len(Args) != 1 {
		usage()
		return
	}

	switch Args[0] {
	case "m", "mem":
		memstat()
	case "c", "cpu":
		cpustat()
	case "d", "disk":
		diskstat()
	case "u":
		uptime()
	default:
		usage()
	}
}

func usage() {
	fmt.Println("\nInvalid Arguments. Choose from below list")
	fmt.Println("\tm - Memory stats")
	fmt.Println("\tc - CPU stats(not yet implemented)")
	fmt.Println("\td - Disk stat")
	fmt.Println("\tu - System uptime\n")
}
