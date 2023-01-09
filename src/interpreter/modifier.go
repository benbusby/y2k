package interpreter

import (
	"github.com/benbusby/y2k/src/utils"
	"math"
	"reflect"
	"strings"
)

type Y2KMod struct {
	VarID    uint8
	ModFn    uint8
	ArgIsVar bool
	ModSize  uint8
	value    string
}

// modMap holds an int->function mapping to match timestamp input
// to the appropriate function to perform on the specified variable.
var modMap = map[uint8]func(*Y2KVar, string, float64){
	1: AddToVar,
	2: SubtractFromVar,
	3: MultiplyVar,
	4: DivideVar,
	5: PowVar,
	9: SetVar,
}

// AddToVar directly modifies a variable by adding a second value to either its
// numVal or strVal property (depending on variable data type).
func AddToVar(y2kVar *Y2KVar, strVal string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal += strVal
		break
	default:
		y2kVar.numVal += numVal
		break
	}
}

// SubtractFromVar modifies a variable by subtracting from the variable's value.
// For strings, this results in a substring from 0:length-N. For all other
// variable types, this is regular subtraction.
func SubtractFromVar(y2kVar *Y2KVar, _ string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = y2kVar.strVal[0 : len(y2kVar.strVal)-int(numVal)]
		break
	default:
		y2kVar.numVal -= numVal
		break
	}
}

// MultiplyVar directly modifies a variable by multiplying the value by a
// number. For strings, this results in a string that is repeated N number of
// times. For all other variable types, this is regular multiplication. Note
// that in this case, val is always treated as a number, even for string
// variables.
func MultiplyVar(y2kVar *Y2KVar, _ string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = strings.Repeat(y2kVar.strVal, int(numVal))
		break
	default:
		y2kVar.numVal *= numVal
		break
	}
}

// DivideVar modifies a variable by dividing the value by a number (if the
// variable is numeric) or a string (if the variable is a string). For strings,
// this results in a string with all instances of the specified string removed.
// For all other variable types, this is regular division.
// value Example: "hello world!" / "o" -> "hell wrld!"
func DivideVar(y2kVar *Y2KVar, strVal string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = strings.ReplaceAll(
			y2kVar.strVal,
			strVal,
			"")
		break
	default:
		y2kVar.numVal /= numVal
		break
	}
}

// PowVar returns the result of exponentiation with a variable's numeric
// value as a base, and numVal input as the exponent.
// This only applies to numeric variables -- string variables are ignored.
func PowVar(y2kVar *Y2KVar, _ string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		return
	default:
		y2kVar.numVal = math.Pow(y2kVar.numVal, numVal)
	}
}

// SetVar overwrites a variable's value with the given input. Note that you
// cannot overwrite a string variable with a numeric value. You would want
// to create a new variable (command 8) with the new data type in that case.
func SetVar(y2kVar *Y2KVar, strVal string, numVal float64) {
	switch y2kVar.Type {
	case Y2KString:
		y2kVar.strVal = strVal
		break
	default:
		y2kVar.numVal = numVal
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
	y2k.DebugMsg("ParseModify: [%s]%s",
		input,
		timestamp[y2k.Digits:],
	)

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

		// Retrieve the possible str and num values of the provided values
		splitValue := utils.SplitStrByN(varMod.value, y2k.Digits)
		strVal := utils.StrArrToPrintable(splitValue)
		numVal := utils.StrArrToFloat(splitValue)

		// If the user specified that the argument is a variable, use the
		// provided input as a variable ID lookup and overwrite the values
		// determined earlier
		if varMod.ArgIsVar {
			argVar := GetVar(uint8(utils.StrArrToInt(splitValue)))
			strVal, numVal = argVar.GetValues()
		}

		modMap[varMod.ModFn](targetVar, strVal, numVal)

		return timestamp
	}

	return y2k.ParseModify(timestamp[y2k.Digits:], reflect.ValueOf(varMod))
}
