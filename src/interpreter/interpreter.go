package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

type Y2K struct {
	Debug     bool
	Digits    int
	Timestamp string
}

// Y2KCommand is an enum to indicate how the interpreter should respond to
// program inputs.
type Y2KCommand uint8

const (
	PRINT  Y2KCommand = 9
	SET    Y2KCommand = 8
	MODIFY Y2KCommand = 7
	WHILE  Y2KCommand = 6
)

// DebugMsg is used for printing useful info about what operations the
// interpreter is performing, and inspecting the values from the timestamps
// that are being interpreted.
func (y2k *Y2K) DebugMsg(prefixSpaces int, msg string) {
	if y2k.Debug {
		fmt.Println(strings.Repeat(" ", prefixSpaces), msg)
	}
}

func (y2k *Y2K) OutputMsg(msg string) {
	debugPrefix := ""
	if y2k.Debug {
		debugPrefix = "    OUTPUT: "
	}

	fmt.Println(fmt.Sprintf("%s%s", debugPrefix, msg))
}

// Parse manages interpreter state and hands off timestamp parsing to the
// appropriate function when changes to interpreter state are made.
// For example, creation of a variable jumps from STANDBY to SET states,
// and moves timestamp parsing to ParseVariable until that function passes
// parsing back to Parse.
func (y2k *Y2K) Parse(timestamp string) {
	// Extract a portion of the timestamp, with size determined by the
	// Y2K.Digits field.
	y2k.DebugMsg(0, fmt.Sprintf("Parse: [%s]%s",
		timestamp[:y2k.Digits],
		timestamp[y2k.Digits:]))
	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	switch Y2KCommand(command) {
	case PRINT:
		y2k.DebugMsg(4, fmt.Sprintf("(%d->ParsePrint)", command))
		timestamp = y2k.ParsePrint(timestamp[y2k.Digits:], Y2KPrint{})
		break
	case SET:
		y2k.DebugMsg(4, fmt.Sprintf("(%d->ParseVariable)", command))
		timestamp = y2k.ParseVariable(timestamp[y2k.Digits:], Y2KVar{})
		break
	case MODIFY:
		y2k.DebugMsg(4, fmt.Sprintf("(%d->ParseModify)", command))
		timestamp = y2k.ParseModify(timestamp[y2k.Digits:], Y2KMod{})
		break
	case WHILE:
		y2k.DebugMsg(4, fmt.Sprintf("(%d->ParseWhile)", command))
		timestamp = y2k.ParseWhile(timestamp[y2k.Digits:], Y2KWhile{})
		break
	}

	if y2k.Digits > len(timestamp)-y2k.Digits {
		// Finished parsing
		return
	}

	y2k.Parse(timestamp[y2k.Digits:])
}
