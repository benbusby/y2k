package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var Y2KExt = ".y2k"
var Printable = " abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"1234567890" +
	"!@#$%^&*()+-<>.,"
var MaxTimestamp = int64(999999999999999999)
var StrTerm = "  "
var LoopTerm = "1999"
var CondTerm = "2000"
var ContinueCmd = "continue"
var DebugDivider = "=============================="

func GetFileModTime(path string, zeroPad bool) string {
	info, err := os.Stat(path)

	if err == nil {
		prefix := ""
		if zeroPad {
			prefix = "0"
		}
		timestamp := info.ModTime().UnixNano()
		if timestamp > MaxTimestamp {
			// File was created after the year 2000, which isn't possible,
			// so let's ignore it
			return ""
		}

		return fmt.Sprintf(prefix+"%d", timestamp)
	}

	return ""
}

func GetCondTerm(loop bool) string {
	if loop {
		return LoopTerm
	}

	return CondTerm
}

func StrToInt(input string) int {
	numVal, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}

	return numVal
}

func StrToFloat(input string) float64 {
	numVal, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	}

	return numVal
}

func FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}

func StrArrToInt(input []string) int {
	numVal, err := strconv.Atoi(strings.Join(input, ""))
	if err != nil {
		return 0
	}

	return numVal
}

func StrArrToFloat(input []string) float64 {
	numVal, err := strconv.ParseFloat(strings.Join(input, ""), 64)
	if err != nil {
		return 0
	}

	return numVal
}

func StrArrToPrintable(input []string) string {
	output := ""
	for _, val := range input {
		index := StrToInt(val)
		if index < len(Printable) {
			output += string(Printable[index])
		}
	}

	return output
}

func SplitStrByN(input string, n int) []string {
	var output []string

	for len(input) != 0 && n < len(input) {
		output = append(output, input[:n])
		input = input[n:]
	}

	output = append(output, input)
	return output
}

func GetFileTimestamp(file string, digits int) string {
	// Check to see if this file is a timestamp-only file (which is the case
	// if GetFileModTime finds a timestamp pre-2000) or if it's a "raw" file
	fileModTime := GetFileModTime(file, digits > 1)

	if len(fileModTime) > 0 {
		return fileModTime
	}

	// File was made after 2000, so we can assume it's likely a raw file
	return ReadY2KRawFile(file)
}

func GetTimestamps(dir string, digits int) string {
	var fullTimestamp = ""
	files, err := os.ReadDir(dir)

	// If the input is not a directory, try reading it as a file
	if err != nil {
		return GetFileTimestamp(dir, digits)
	}

	directoryPath, _ := filepath.Abs(dir)

	// Sort contents of the specified directory by name.
	// Y2K files should be named in an easily sortable manner when creating
	// programs (i.e. 00.y2k -> 01.y2k -> etc).
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		// Ignore any non *.y2k files
		if !strings.HasSuffix(file.Name(), Y2KExt) {
			continue
		}

		// Append timestamp to slice
		fullPath := filepath.Join(directoryPath, file.Name())
		timestamp := GetFileModTime(fullPath, digits > 1)
		if len(fullTimestamp) != 0 && len(timestamp) > 0 {
			// Snip off the leading digit for all timestamps except
			// the first one. We do this to avoid issues with commands
			// spanning across multiple files, where the next desired
			// digit might be a "0" (which would be ignored in a timestamp)
			timestamp = timestamp[digits:]
		}
		fullTimestamp += timestamp
	}

	return fullTimestamp
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
