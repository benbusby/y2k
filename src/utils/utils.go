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
var Printable = " abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()"
var PrintStringTerm = "  "

func GetFileModTime(path string, zeroPad bool) string {
	info, err := os.Stat(path)

	if err == nil {
		prefix := ""
		if zeroPad {
			prefix = "0"
		}
		return fmt.Sprintf(prefix+"%d", info.ModTime().UnixNano())
	}

	return "00"
}

func SafeStrToInt(input string) int {
	numVal, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}

	return numVal
}

func SafeStrArrToInt(input []string) int {
	numVal, err := strconv.Atoi(strings.Join(input, ""))
	if err != nil {
		return 0
	}

	return numVal
}

func StrArrToPrintable(input []string) string {
	output := ""
	for _, val := range input {
		output += string(Printable[SafeStrToInt(val)])
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

func GetDirTimestamps(dir string, zeroPad bool) []string {
	var timestamps []string
	files, err := os.ReadDir(dir)

	// Ensure the directory is actually valid before continuing
	if err != nil {
		panic(fmt.Sprintf("\"%s\" is not a valid directory", dir))
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
		timestamp := GetFileModTime(fullPath, zeroPad)
		timestamps = append(timestamps, timestamp)
	}

	return timestamps
}
