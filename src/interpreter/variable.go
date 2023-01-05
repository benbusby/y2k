package interpreter

import (
	"github.com/benbusby/y2k/src/utils"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var VarMap = map[uint8]*Y2KVar{}

// Y2KVarType is an enum to indicate how the interpreter should treat a Y2KVar.
type Y2KVarType uint8

const (
	Y2KString  Y2KVarType = 1
	Y2KInt     Y2KVarType = 2
	Y2KFloat   Y2KVarType = 3
	Y2KVarCopy Y2KVarType = 9
)

// Y2KVar is a struct for all variables created by Y2K programs. These contain
// both numeric and string values as well as a data type. When creating numeric
// variables, the strVal property is used to construct a numeric value while
// parsing, until the variable's Size is reached.
type Y2KVar struct {
	ID     uint8
	Type   Y2KVarType
	Size   uint8
	strVal string
	numVal float64
}

// GetValue returns the appropriate value for a particular variable. If it's a
// numeric variable, it returns the numeric value, otherwise it returns the
// string value.
func (y2kVar *Y2KVar) GetValue() string {
	switch y2kVar.Type {
	case Y2KString:
		return y2kVar.strVal
	default:
		return utils.FloatToString(y2kVar.numVal)
	}
}

// GetValues returns both strVal and numVal of a variable.
func (y2kVar *Y2KVar) GetValues() (string, float64) {
	return y2kVar.strVal, y2kVar.numVal
}

// GetVar retrieves a variable from the existing ID->var map,
// or returns an empty version of the variable struct if the
// request var id has not been set.
func GetVar(id uint8) *Y2KVar {
	if variable, ok := VarMap[id]; ok {
		return variable
	}

	// If the variable has not been set yet, insert it now.
	VarMap[id] = &Y2KVar{}
	return VarMap[id]
}

// FromCLIArg takes a command line argument and turns it into a variable for the
// programs to reference as needed. Variables added from the command line are
// inserted into the map backwards from the map's max index (9 for 1-digit
// parsing, 99 for 2-digit parsing, etc).
func (y2k Y2K) FromCLIArg(input string) {
	// Determine if the argument is a string or numeric.
	// Assume the variable is numeric, unless a non-numeric other than '.' is
	// found.
	argType := Y2KInt
	for _, c := range input {
		if unicode.IsLetter(c) && c != '.' {
			argType = Y2KString
		}
	}

	// Command line variables are added to the end of the map, which depends on
	// the number of digits that are parsed at one time (a parsing size of 1
	// should insert variables from 9->8->etc, a parsing size of 2 should insert
	// from 99->98->etc.)
	mapInd, _ := strconv.Atoi(strings.Repeat("9", y2k.Digits))
	for VarMap[uint8(mapInd)] != nil {
		mapInd -= 1
	}

	// Finalize and insert the new var into the previously determined index
	VarMap[uint8(mapInd)] = &Y2KVar{
		ID:     uint8(mapInd),
		Size:   uint8(len(input)),
		strVal: input,
		numVal: utils.StrToFloat(input),
		Type:   argType,
	}
}

// ParseVariable recursively builds a new Y2KVar to insert into the global
// variable map.
// The variable creation process follows a specific order:
//
// start creation -> set ID -> set type -> set size -> read values
//
// So to create a numeric variable with the value 100 and an ID of 1, the
// chain of values would need to be:
//
// 3 1 2 3 1 0 0
func (y2k Y2K) ParseVariable(timestamp string, val reflect.Value) string {
	newVar := val.Interface().(Y2KVar)
	input := timestamp[:y2k.Digits]

	y2k.DebugMsg("ParseVariable: [%s]%s",
		input,
		timestamp[y2k.Digits:],
	)

	// Regardless of data type, var values are created as a string first, in
	// order to sequentially create the variable value across multiple passes
	// of the parser (i.e. 100 has to be split between multiple passes, so "1"
	// is added first, then "0", then the last "0", then converted to an
	// integer).
	if newVar.Type == Y2KString {
		input = string(utils.Printable[utils.StrToInt(input)])
	}
	newVar.strVal += input

	if len(newVar.strVal) >= int(newVar.Size) {
		newVar.strVal = newVar.strVal[:newVar.Size]

		if newVar.Type == Y2KVarCopy {
			copyVar := GetVar(uint8(utils.StrToInt(newVar.strVal)))
			newVar.Type = copyVar.Type
			newVar.Size = copyVar.Size
			newVar.numVal = copyVar.numVal
			newVar.strVal = copyVar.strVal
		} else {
			// Init numeric value of variable
			if newVar.Type == Y2KFloat {
				// First digit of a float is where the decimal should be placed
				decimalIndex := utils.StrToInt(newVar.strVal[0:1])
				newVar.strVal = newVar.strVal[1:decimalIndex+1] +
					"." +
					newVar.strVal[decimalIndex+1:]
			}

			newVar.numVal = utils.StrToFloat(newVar.strVal)
		}

		// Insert finished variable into variable map
		VarMap[newVar.ID] = &newVar

		// Return handling of the parser back to Parse
		return timestamp
	}

	return y2k.ParseVariable(timestamp[y2k.Digits:], reflect.ValueOf(newVar))
}
