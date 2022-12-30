package main

import (
	"flag"
	"fmt"
	"github.com/benbusby/y2k/src/interpreter"
	"github.com/benbusby/y2k/src/utils"
)

func main() {
	var timestamp string
	digits := flag.Int(
		"d",
		1,
		"Set # of digits to parse at a time")
	debug := flag.Bool(
		"debug",
		false,
		"Enable to view interpreter steps in console")
	flag.Parse()

	y2k := &interpreter.Y2K{Digits: *digits, Debug: *debug}

	for _, arg := range flag.Args() {
		// Assume first argument is the directory to use for parsing
		if len(timestamp) == 0 {
			timestamp = utils.GetDirTimestamps(arg, *digits)
			continue
		}

		y2k.FromCLIArg(arg)
	}

	if len(timestamp) == 0 {
		fmt.Println("Missing input dir!\n\nUsage: y2k <directory> [args]")
		flag.PrintDefaults()
		return
	}

	y2k.Parse(timestamp)
}
