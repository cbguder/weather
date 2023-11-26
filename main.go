package main

import (
	"os"

	"github.com/cbguder/weather/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
