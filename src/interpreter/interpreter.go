package interpreter

import (
	"fmt"
	"reflect"
	"strings"
	"y2k/src/utils"
)

type Y2K struct {
	Debug  bool
	Digits int
}

type Instruction struct {
	val reflect.Value
	fn  func(Y2K, string, reflect.Value) string
}

type Y2KCommand uint8

const (
	PRINT  Y2KCommand = 9
	SET    Y2KCommand = 8
	MODIFY Y2KCommand = 7
	WHILE  Y2KCommand = 6
)

var instMap map[Y2KCommand]Instruction

// CreateStruct uses reflection to form a struct from N-sized chunks
// from the timestamp. The struct that is constructed is mapped to
// a Y2KCommand and holds all values that are relevant to performing
// the specified command (i.e. Y2KVar establishes variable ID,
// size, and type).
func (y2k Y2K) CreateStruct(
	timestamp string,
	v reflect.Value,
) (reflect.Value, string) {
	modFields := 0

	for i := 0; i < v.NumField(); i++ {
		// Ignore private struct fields
		if !v.Field(i).CanSet() {
			continue
		}

		idx := y2k.Digits * modFields
		val := utils.StrToInt(timestamp[idx : idx+y2k.Digits])

		switch v.Field(i).Type().Kind() {
		case reflect.Int:
		case reflect.Int8:
			v.Field(i).SetInt(int64(val))
			break
		case reflect.Uint:
		case reflect.Uint8:
			v.Field(i).SetUint(uint64(val))
			break
		default:
			panic(fmt.Sprintf(
				"Unhandled type reflection: %s (in %s)",
				v.Field(i).Type().Kind(),
				v.String()))
		}

		modFields += 1
	}

	newStart := y2k.Digits * modFields
	return v, timestamp[newStart:]
}

// DebugMsg is used for printing useful info about what operations the
// interpreter is performing, and inspecting the values from the timestamps
// that are being interpreted.
func (y2k Y2K) DebugMsg(prefixSpaces int, msg string) {
	if y2k.Debug {
		fmt.Println(strings.Repeat(" ", prefixSpaces), msg)
	}
}

func (y2k Y2K) OutputMsg(msg string) {
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
func (y2k Y2K) Parse(timestamp string) {
	// Extract a portion of the timestamp, with size determined by the
	// Y2K.Digits field.
	y2k.DebugMsg(0, fmt.Sprintf("Parse: [%s]%s",
		timestamp[:y2k.Digits],
		timestamp[y2k.Digits:]))
	command := Y2KCommand(utils.StrToInt(timestamp[:y2k.Digits]))

	if instruction, ok := instMap[command]; ok {
		var y2kStruct reflect.Value
		y2kStruct, timestamp = y2k.CreateStruct(
			timestamp[y2k.Digits:],
			instruction.val)
		timestamp = instruction.fn(y2k, timestamp, y2kStruct)
	}

	if y2k.Digits > len(timestamp)-y2k.Digits {
		// Finished parsing
		return
	}

	y2k.Parse(timestamp[y2k.Digits:])
}

func init() {
	instMap = map[Y2KCommand]Instruction{
		PRINT:  {reflect.ValueOf(&Y2KPrint{}).Elem(), Y2K.ParsePrint},
		SET:    {reflect.ValueOf(&Y2KVar{}).Elem(), Y2K.ParseVariable},
		MODIFY: {reflect.ValueOf(&Y2KMod{}).Elem(), Y2K.ParseModify},
		WHILE:  {reflect.ValueOf(&Y2KWhile{}).Elem(), Y2K.ParseWhile},
	}
}
