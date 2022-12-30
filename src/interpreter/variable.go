package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"y2k/src/utils"
)

var VarMap = map[uint8]*Y2KVar{}

// Y2KVarType is an enum to indicate how the interpreter should treat a Y2KVar.
type Y2KVarType uint8

const (
	Y2KString  Y2KVarType = 1
	Y2KNumber  Y2KVarType = 2
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
	intVal int
	strVal string
}

// GetValue returns the appropriate value for a particular variable. If it's a
// numeric variable, it returns the numeric value, otherwise it returns the
// string value.
func (y2kVar *Y2KVar) GetValue() string {
	switch y2kVar.Type {
	case Y2KString:
		return y2kVar.strVal
	case Y2KNumber:
		return strconv.Itoa(y2kVar.intVal)
	}

	return ""
}

// GetVar retrieves a variable from the existing ID->var map,
// or returns an empty version of the variable struct if the
// request var id has not been set.
func GetVar(id uint8) *Y2KVar {
	if variable, ok := VarMap[id]; ok {
		return variable
	}

	// Alert the user if the variable they're attempting to use has not been set yet.
	warnMsg := fmt.Sprintf("WARNING -- Variable [%d] does not exist!", id)
	if id == 9 {
		warnMsg += "\n(probably a command line argument...)"
	}

	fmt.Println(warnMsg)
	return &Y2KVar{}
}

// FromCLIArg takes a command line argument and turns it into a variable for the
// programs to reference as needed. Variables added from the command line are
// inserted into the map backwards from the map's max index (9 for 1-digit
// parsing, 99 for 2-digit parsing, etc).
func (y2k Y2K) FromCLIArg(input string) {
	// Determine if the argument is a string or numeric.
	// Assume the variable is numeric, unless a non-numeric other than '.' is
	// found.
	argType := Y2KNumber
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
		intVal: utils.StrToInt(input),
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

	y2k.DebugMsg(fmt.Sprintf("ParseVariable: [%s]%s",
		input,
		timestamp[y2k.Digits:]))

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
		if newVar.Type == Y2KVarCopy {
			copyVar := GetVar(uint8(utils.StrToInt(newVar.strVal)))
			newVar.Type = copyVar.Type
			newVar.Size = copyVar.Size
			newVar.intVal = copyVar.intVal
			newVar.strVal = copyVar.strVal
		} else {
			// Init int value of variable
			newVar.intVal = utils.StrToInt(newVar.strVal[:newVar.Size])
		}

		// Insert finished variable into variable map
		VarMap[newVar.ID] = &newVar

		// Return handling of the parser back to Parse
		return timestamp
	}

	return y2k.ParseVariable(timestamp[y2k.Digits:], reflect.ValueOf(newVar))
}
