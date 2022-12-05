package interpreter

import (
	"fmt"
	"strconv"
)

type Y2K struct {
	Debug      bool
	Digits     int
	Timestamps []string
}

// Y2KCommand is an enum to indicate how the interpreter should respond to
// program inputs.
type Y2KCommand uint8

const (
	PRINT  Y2KCommand = 9
	SET    Y2KCommand = 8
	MODIFY Y2KCommand = 7
)

// DebugMsg is used for printing useful info about what operations the
// interpreter is performing, and inspecting the values from the timestamps
// that are being interpreted.
func (y2k *Y2K) DebugMsg(msg string) {
	if y2k.Debug {
		fmt.Println(msg)
	}
}

// Parse manages interpreter state and hands off timestamp parsing to the
// appropriate function when changes to interpreter state are made.
// For example, creation of a variable jumps from STANDBY to SET states,
// and moves timestamp parsing to ParseVariable until that function passes
// parsing back to Parse.
func (y2k *Y2K) Parse(timestamp string) {
	if y2k.Digits > len(timestamp) {
		// Finished parsing
		return
	}

	// Extract a portion of the timestamp, with size determined by the Y2K.Digits field.
	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	switch Y2KCommand(command) {
	case PRINT:
		y2k.DebugMsg(fmt.Sprintf("%d: Print", command))
		timestamp = y2k.ParsePrint(timestamp[y2k.Digits:], Y2KPrint{})
		break
	case SET:
		y2k.DebugMsg(fmt.Sprintf("%d: Create Variable", command))
		timestamp = y2k.ParseVariable(timestamp[y2k.Digits:], Y2KVar{})
		break
	case MODIFY:
		y2k.DebugMsg(fmt.Sprintf("%d: Modify Variable", command))
		//timestamp = y2k.ParseModify(timestamp[y2k.Digits:], Y2KMod{})
		break
	}

	y2k.Parse(timestamp[y2k.Digits:])
}

func (y2k *Y2K) Run() {
	// FIXME: Refactor this to use one timestamp created from multiple files,
	// rather than a list of timestamps. This will eliminate the issues around
	// commands that span multiple files.
	for _, timestamp := range y2k.Timestamps {
		y2k.Parse(timestamp)
	}
}
