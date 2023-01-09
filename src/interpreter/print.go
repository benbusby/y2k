package interpreter

import (
	"github.com/benbusby/y2k/src/utils"
	"reflect"
)

// Y2KPrintType is an enum to indicate to the interpreter what should be printed.
type Y2KPrintType uint8

const (
	Y2KPrintString Y2KPrintType = 1
	Y2KPrintVar    Y2KPrintType = 2
)

type Y2KPrint struct {
	Type  Y2KPrintType
	Size  int
	value string
}

func (y2k Y2K) ParsePrint(timestamp string, val reflect.Value) string {
	y2kPrint := val.Interface().(Y2KPrint)

	input := timestamp[:y2k.Digits]
	y2k.DebugMsg("ParsePrint: [%s]%s",
		input,
		timestamp[y2k.Digits:],
	)

	y2kPrint.value += input

	if len(y2kPrint.value) >= y2kPrint.Size*y2k.Digits {
		// If we're printing a variable, the value will be an integer
		// variable ID to print. Otherwise, we need to split the string
		// into N-sized chunks (dependent on interpreter parsing window
		// size) and print each character that matches each digit.
		switch y2kPrint.Type {
		case Y2KPrintString:
			splitValues := utils.SplitStrByN(y2kPrint.value, y2k.Digits)
			strValue := utils.StrArrToPrintable(splitValues)
			y2k.OutputMsg(strValue)
			break
		case Y2KPrintVar:
			printVar := GetVar(uint8(utils.StrToInt(y2kPrint.value)))
			y2k.OutputMsg(printVar.GetValue())
			break
		}

		return timestamp
	}

	return y2k.ParsePrint(timestamp[y2k.Digits:], reflect.ValueOf(y2kPrint))
}
