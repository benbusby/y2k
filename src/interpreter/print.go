package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"y2k/src/utils"
)

// Y2KPrintType is an enum to indicate to the interpreter what should be printed.
type Y2KPrintType uint8

const (
	Y2KPrintNone   Y2KPrintType = 0
	Y2KPrintString Y2KPrintType = 1
	Y2KPrintVar    Y2KPrintType = 2
)

type Y2KPrint struct {
	Type   Y2KPrintType
	String string
}

func (y2k *Y2K) ParsePrint(timestamp string, print Y2KPrint) string {
	// If a print statement aligns with the end of the timestamp, we can at least
	// just print the contents of Y2KPrint.String, if the user was wanting to print
	// a string.
	if y2k.Digits > len(timestamp) {
		if print.Type == Y2KPrintString {
			fmt.Println(print.String)
		}
		return timestamp
	}

	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	if print.Type == Y2KPrintNone {
		// Guard against invalid print type
		if command == 0 {
			panic("Cannot set print type to 0")
		}

		print.Type = Y2KPrintType(command)
		y2k.DebugMsg(fmt.Sprintf("  Set Print Type: %d", print.Type))
	} else {
		// If we're printing a variable, the next input will be the variable ID.
		// We can use that to print the variable value and return the timestamp
		// back to the caller.
		if print.Type == Y2KPrintVar {
			y2k.DebugMsg(fmt.Sprintf("  Print Var: %d", command))
			variable := VarMap[uint8(command)]
			fmt.Println(variable.GetValue())

			return timestamp
		}

		// Otherwise we need to begin building a string until there have been two
		// back-to-back spaces (two 0 inputs). This is just an arbitrary way of
		// determining when parsing of a print string should end.
		print.String += string(utils.Printable[command])

		if strings.HasSuffix(print.String, utils.PrintStringTerm) {
			fmt.Println(print.String[0 : len(print.String)-len(utils.PrintStringTerm)])
			return timestamp
		}
	}

	return y2k.ParsePrint(timestamp[y2k.Digits:], print)
}