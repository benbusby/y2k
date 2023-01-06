package interpreter

import (
	"fmt"
	"github.com/benbusby/y2k/src/utils"
	"reflect"
	"strconv"
	"strings"
)

// Y2KPrintType is an enum to indicate to the interpreter what should be printed.
type Y2KPrintType uint8

const (
	Y2KPrintString Y2KPrintType = 1
	Y2KPrintVar    Y2KPrintType = 2
)

type Y2KPrint struct {
	Type Y2KPrintType
	out  string
}

func (y2k Y2K) ParsePrint(timestamp string, val reflect.Value) string {
	y2kPrint := val.Interface().(Y2KPrint)

	input := utils.StrToInt(timestamp[:y2k.Digits])
	y2k.DebugMsg("ParsePrint: [%s]%s",
		timestamp[:y2k.Digits],
		timestamp[y2k.Digits:],
	)

	// If we're printing a variable, the next input will be the variable ID.
	// We can use that to print the variable value and return the timestamp
	// back to the caller.
	if y2kPrint.Type == Y2KPrintVar {
		y2k.DebugMsg("Print Var[%s]", strconv.Itoa(input))
		variable := GetVar(uint8(input))
		y2k.OutputMsg(variable.GetValue())

		return timestamp
	} else if y2kPrint.Type != Y2KPrintString {
		panic(fmt.Sprintf("Unknown print type: %d", y2kPrint.Type))
	}

	// Otherwise we need to begin building a string until there have been two
	// back-to-back spaces (two 0 inputs). This is just an arbitrary way of
	// determining when parsing of a print string should end.
	y2kPrint.out += string(utils.Printable[input])

	// Check if the string terminator (whitespace * N-digits) has been added,
	// and if so, strip out the terminator and print the string.
	if strings.HasSuffix(y2kPrint.out, utils.StrTerm) {
		y2k.OutputMsg(y2kPrint.out[0 : len(y2kPrint.out)-len(utils.StrTerm)])
		return timestamp
	}

	return y2k.ParsePrint(timestamp[y2k.Digits:], reflect.ValueOf(y2kPrint))
}
