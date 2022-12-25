package interpreter

import (
	"reflect"
	"strings"
	"y2k/src/utils"
)

type Y2KMod struct {
	VarID   uint8
	ModFn   uint8
	ModSize uint8
	value   string
}

// modMap holds an int->function mapping to match timestamp input
// to the appropriate function to perform on the specified variable.
var modMap = map[uint8]func(*Y2KVar, []string){
	1: AddToVar,
	2: SubtractFromVar,
	3: MultiplyVar,
	4: DivideVar,
	5: AddVarToVar,
	6: CopyFromVar,
}

// AddToVar directly modifies a variable by adding a second value to either its
// intVal or strVal property (depending on variable data type).
func AddToVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal += utils.StrArrToPrintable(values)
		break
	default:
		y2kVar.intVal += utils.StrArrToInt(values)
		break
	}
}

func AddVarToVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal += GetVar(uint8(utils.StrArrToInt(values))).strVal
		break
	default:
		y2kVar.intVal += GetVar(uint8(utils.StrArrToInt(values))).intVal
		break
	}
}

func CopyFromVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = GetVar(uint8(utils.StrArrToInt(values))).strVal
		break
	default:
		y2kVar.intVal = GetVar(uint8(utils.StrArrToInt(values))).intVal
		break
	}
}

// SubtractFromVar modifies a variable by subtracting from the variable's value.
// For strings, this results in a substring from 0:length-N. For all other
// variable types, this is regular subtraction.
func SubtractFromVar(y2kVar *Y2KVar, values []string) {
	intVal := utils.StrArrToInt(values)
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = y2kVar.strVal[0 : len(y2kVar.strVal)-intVal]
		break
	default:
		y2kVar.intVal -= intVal
		break
	}
}

// MultiplyVar directly modifies a variable by multiplying the value by a
// number. For strings, this results in a string that is repeated N number of
// times. For all other variable types, this is regular multiplication. Note
// that in this case, val is always treated as a number, even for string
// variables.
func MultiplyVar(y2kVar *Y2KVar, values []string) {
	intVal := utils.StrArrToInt(values)

	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = strings.Repeat(y2kVar.strVal, intVal)
		break
	default:
		y2kVar.intVal *= intVal
		break
	}
}

// DivideVar modifies a variable by dividing the value by a number (if the
// variable is numeric) or a string (if the variable is a string). For strings,
// this results in a string with all instances of the specified string removed.
// For all other variable types, this is regular division.
// out Example: "hello world!" / "o" -> "hell wrld!"
func DivideVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = strings.ReplaceAll(
			y2kVar.strVal,
			utils.StrArrToPrintable(values),
			"")
		break
	default:
		y2kVar.intVal /= utils.StrArrToInt(values)
		break
	}
}

// ParseModify recursively builds a set of values to modify an existing
// variable. The order of values are:
//
//	<target variable ID> -> <function ID> -> <mod size> -> <mod value>
//
// Once the mod size has been reached, we can pass the mod value to the desired
// function and return the timestamp back to the original caller.
func (y2k Y2K) ParseModify(timestamp string, val reflect.Value) string {
	varMod := val.Interface().(Y2KMod)

	input := timestamp[:y2k.Digits]

	varMod.value += input

	if len(varMod.value) >= int(varMod.ModSize) {
		// Although we have the desired size of the modification, we don't
		// know how the modification value needs to be interpreted. By
		// converting the mod value to a slice of strings, we can pass off
		// final interpretation of the value to the actual function that is
		// performing the modification. For example, adding to a string
		// should interpret inputs as a string ("h" + 9 == "hi"), but
		// multiplying a string should interpret the input as a number.
		// ("h" * 9 == "hhhhhhhhh").
		targetVar := GetVar(varMod.VarID)
		varMod.value = varMod.value[:varMod.ModSize]
		modMap[varMod.ModFn](
			targetVar,
			utils.SplitStrByN(varMod.value, y2k.Digits))

		return timestamp
	}

	return y2k.ParseModify(timestamp[y2k.Digits:], reflect.ValueOf(varMod))
}
