package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var commentChar = "#"

// ReadY2KRawFile reads file contents of a Y2K file instead of file timestamps.
// The syntax for these files should typically follow the following structure:
// <numeric command values> : <comment>
//
// For example, to print the letter "a", you could write:
// 921 : Print "a"
//
// The file is read line by line, and whitespace and comments are ignored, so
// Y2K programs can take up as much space as needed to make sense without
// impacting the interpreter.
func ReadY2KRawFile(file string) string {
	timestamp := ""
	raw, err := os.Open(file)
	check(err)

	defer func(raw *os.File) {
		err := raw.Close()
		check(err)
	}(raw)

	// Strip all whitespace and comments from file
	scanner := bufio.NewScanner(raw)
	for scanner.Scan() {
		line := scanner.Text()

		// Remove any comments from line
		commentIndex := strings.Index(line, commentChar)
		if commentIndex >= 0 {
			line = line[:commentIndex]
		}

		// Remove extra whitespace
		line = strings.ReplaceAll(line, " ", "")

		// Append to timestamp
		timestamp += line
	}

	return timestamp
}

// WriteFileTimestamp creates an empty file at <path>/<fileNum>.y2k and modifies
// the file's timestamp with the value provided.
func WriteFileTimestamp(timestamp string, path string, fileNum int) {
	filename := fmt.Sprintf("%s/%d.y2k", path, fileNum)
	file, err := os.Create(filename)
	check(err)

	err = file.Close()
	check(err)

	// Prepend a digit for all file timestamps after the first file. The reason
	// for this is explained in the README.
	if fileNum > 0 {
		timestamp = "8" + timestamp
	}

	if len(timestamp) != 18 {
		panic("Error: Invalid timestamp length -- must be 18 chars long")
	}

	fileTime := time.Unix(int64(StrToInt(timestamp[:9])), int64(StrToInt(timestamp[9:])))

	fmt.Println(fmt.Sprintf("Writing %s -- %s (%s)", filename, timestamp, fileTime))

	err = os.Chtimes(filename, fileTime, fileTime)
	check(err)
}

// ExportRawToTimestampFiles takes the timestamp created from a raw Y2K file
// and outputs a set of empty files that have their timestamps modified to
// perform the same operations as the raw file.
func ExportRawToTimestampFiles(timestamp string, path string) {
	files := 0

	// Ensure path exists, and create it if not
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		check(err)
	}

	for len(timestamp) > 0 {
		maxLen := 17
		if files == 0 {
			maxLen = 18
		}

		// Ensure the timestamp has trailing 0s (not leading, which would
		// impact multi-file commands) if it's shorter than the maximum
		// length. This typically happens when programs require only part
		// of an additional file's timestamp to work properly.
		for len(timestamp) < maxLen {
			timestamp += "0"
		}

		WriteFileTimestamp(timestamp[:maxLen], path, files)
		timestamp = timestamp[maxLen:]
		files += 1
	}
}
