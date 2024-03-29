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
	export := flag.Bool(
		"export",
		false,
		"Export a Y2K raw file to a set of timestamp-only files")
	outdir := flag.String(
		"outdir",
		"./y2k-out",
		"Set the output directory for timestamp-only files when exporting a raw Y2K file.\n"+
			"This directory will be created if it does not exist.")
	flag.Parse()

	y2k := &interpreter.Y2K{Digits: *digits, Debug: *debug}

	for _, arg := range flag.Args() {
		// Assume first argument is the directory or file to use for parsing
		if len(timestamp) == 0 {
			if *export {
				// If we're exporting, assume we're only reading raw Y2K file
				// contents, and export to a set of empty files.
				timestamp = utils.ReadY2KRawFile(arg)

				utils.ExportRawToTimestampFiles(timestamp, *outdir)
				return
			} else {
				timestamp = utils.GetTimestamps(arg, *digits)
			}
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
