package main

import (
	"flag"
	"fmt"
	"y2k/src/interpreter"
	"y2k/src/utils"
)

func main() {
	digits := flag.Int(
		"d",
		2,
		"Set # of digits to parse at a time")
	debug := flag.Bool(
		"debug",
		false,
		"Enable to view interpreter steps in console")
	noTrim := flag.Bool(
		"no-trim",
		false,
		"Disables trimming of the first N digits "+
			"of file timestamps after the first file")
	flag.Parse()

	y2k := &interpreter.Y2K{Digits: *digits, Debug: *debug}

	for _, arg := range flag.Args() {
		// Assume first argument is the directory to use for parsing
		if len(y2k.Timestamp) == 0 {
			y2k.Timestamp = utils.GetDirTimestamps(arg, *digits, *noTrim)
			continue
		}

		interpreter.FromCLIArg(arg, *digits)
	}

	if len(y2k.Timestamp) == 0 {
		fmt.Println("Missing input dir!\n\nUsage: y2k <directory> [args]")
		flag.PrintDefaults()
		return
	}

	y2k.Parse(y2k.Timestamp)
}
