package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"y2k/src/utils"
)

type Y2KMod struct {
	VarID       uint8
	ModSize     uint8
	ModValue    string
	ModFunction func(*Y2KVar, []string)
}

// modMap holds an int->function mapping to match timestamp input
// to the appropriate function to perform on the specified variable.
var modMap = map[uint8]func(*Y2KVar, []string){
	1: AddToVar,
	2: SubtractFromVar,
	3: MultiplyVar,
	4: DivideVar,
	5: AddVarToVar,
}

// AddToVar directly modifies a variable by adding a second value to either its
// NumberVal or StringVal property (depending on variable data type).
func AddToVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.StringVal += utils.StrArrToPrintable(values)
		break
	default:
		y2kVar.NumberVal += utils.StrArrToInt(values)
		break
	}
}

func AddVarToVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.StringVal += VarMap[uint8(utils.StrArrToInt(values))].StringVal
		break
	default:
		y2kVar.NumberVal += VarMap[uint8(utils.StrArrToInt(values))].NumberVal
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
		y2kVar.StringVal = y2kVar.StringVal[0 : len(y2kVar.StringVal)-intVal]
		break
	default:
		y2kVar.NumberVal -= intVal
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
		y2kVar.StringVal = strings.Repeat(y2kVar.StringVal, intVal)
		break
	default:
		y2kVar.NumberVal *= intVal
		break
	}
}

// DivideVar modifies a variable by dividing the value by a number (if the
// variable is numeric) or a string (if the variable is a string). For strings,
// this results in a string with all instances of the specified string removed.
// For all other variable types, this is regular division.
// Str Example: "hello world!" / "o" -> "hell wrld!"
func DivideVar(y2kVar *Y2KVar, values []string) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.StringVal = strings.ReplaceAll(
			y2kVar.StringVal,
			utils.StrArrToPrintable(values),
			"")
		break
	default:
		y2kVar.NumberVal /= utils.StrArrToInt(values)
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
func (y2k *Y2K) ParseModify(timestamp string, varMod Y2KMod) string {
	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	if varMod.VarID == 0 {
		varMod.VarID = uint8(command)
		y2k.DebugMsg(4, fmt.Sprintf("Variable ID: %d", varMod.VarID))
	} else if varMod.ModFunction == nil {
		varMod.ModFunction = modMap[uint8(command)]
		y2k.DebugMsg(4, fmt.Sprintf("Function: %d", command))
	} else if varMod.ModSize == 0 {
		varMod.ModSize = uint8(command)
		y2k.DebugMsg(4, fmt.Sprintf("Modifier Size: %d", varMod.ModSize))
	} else {
		strVal := strconv.Itoa(command)
		varMod.ModValue += strVal
		y2k.DebugMsg(4, fmt.Sprintf("(+ value: %s)", strVal))

		if len(varMod.ModValue) >= int(varMod.ModSize) {
			// Although we have the desired size of the modification, we don't
			// know how the modification value needs to be interpreted. By
			// converting the mod value to a slice of strings, we can pass off
			// final interpretation of the value to the actual function that is
			// performing the modification. For example, adding to a string
			// should interpret inputs as a string ("h" + 9 == "hi"), but
			// multiplying a string should interpret the inputs as a number.
			// ("h" * 9 == "hhhhhhhhh").
			targetVar := VarMap[varMod.VarID]
			varMod.ModValue = varMod.ModValue[:varMod.ModSize]
			varMod.ModFunction(
				targetVar,
				utils.SplitStrByN(varMod.ModValue, y2k.Digits))

			return timestamp
		}
	}
	return y2k.ParseModify(timestamp[y2k.Digits:], varMod)
}
