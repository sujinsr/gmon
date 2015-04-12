package main

import (
	"fmt"
	"os"
)

func main() {
	/* get command line arguments without program name */
	Args := os.Args[1:]

	if len(Args) != 1 {
		fmt.Println("Invalid Arguments")
		return
	}

	switch Args[0] {
	case "m", "mem":
		memstat()
	case "c", "cpu":
		cpustat()
	default:
		fmt.Println("Invalid Arguments")
	}
}
