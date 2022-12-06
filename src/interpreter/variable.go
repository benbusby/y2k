package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"y2k/src/utils"
)

var VarMap = map[uint8]*Y2KVar{}

// Y2KVarType is an enum to indicate how the interpreter should treat a Y2KVar.
type Y2KVarType uint8

const (
	Y2KNoType Y2KVarType = 0
	Y2KString Y2KVarType = 1
	Y2KNumber Y2KVarType = 2
)

// Y2KVar is a struct for all variables created by Y2K programs. These contain
// both numeric and string values as well as a data type. When creating numeric
// variables, the StringVal property is used to construct a numeric value while
// parsing, until the variable's Size is reached.
type Y2KVar struct {
	ID        uint8
	Size      uint8
	NumberVal int
	StringVal string
	Type      Y2KVarType
}

// GetValue returns the appropriate value for a particular variable. If it's a
// numeric variable, it returns the numeric value, otherwise it returns the
// string value.
func (y2kVar *Y2KVar) GetValue() string {
	switch y2kVar.Type {
	case Y2KString:
		return y2kVar.StringVal
	case Y2KNumber:
		return strconv.Itoa(y2kVar.NumberVal)
	}

	return ""
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
func (y2k *Y2K) ParseVariable(timestamp string, newVar Y2KVar) string {
	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	if newVar.ID == 0 {
		// Step 1 -- Set variable id
		newVar.ID = uint8(command)
		y2k.DebugMsg(fmt.Sprintf("  Set ID: %d", newVar.ID))
	} else if newVar.Type == Y2KNoType {
		// Step 2 -- Set variable type
		newVar.Type = Y2KVarType(command)
		y2k.DebugMsg(fmt.Sprintf("  Set Type: %d", newVar.Type))
	} else if newVar.Size == 0 {
		// Step 3 -- Set variable size
		newVar.Size = uint8(command)
		y2k.DebugMsg(fmt.Sprintf("  Set Size: %d", newVar.Size))
	} else {
		// Step 4 -- Init new var values
		// Regardless of data type, this is created as a string first, in
		// order to sequentially create the variable value across multiple passes
		// of the parser (i.e. 100 has to be split between multiple passes, so
		// "1" is added first, then "0", then the last "0", then converted to
		// an integer).
		strVal := strconv.Itoa(command)
		newVar.StringVal += strVal
		y2k.DebugMsg(fmt.Sprintf("    (+ value: %s)", strVal))
		if len(newVar.StringVal) >= int(newVar.Size) {
			numericVal, _ := strconv.Atoi(newVar.StringVal[:newVar.Size])
			newVar.NumberVal = numericVal

			y2k.DebugMsg(fmt.Sprintf("  Set Value: %d", newVar.NumberVal))

			// Insert finished variable into variable map
			VarMap[newVar.ID] = &newVar

			// Return handling of the parser back to Parse
			return timestamp
		}
	}

	return y2k.ParseVariable(timestamp[y2k.Digits:], newVar)
}

// FromCLIArg takes a command line argument and turns it into a variable for the
// programs to reference as needed. Variables added from the command line are
// inserted into the map backwards from the map's max index (9 for 1-digit parsing,
// 99 for 2-digit parsing, etc).
func FromCLIArg(input string, parsingSize int) {
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
	// the number of digits that are parsed at one time (a parsing size of 1 should
	// insert variables from 9->8->etc, whereas a parsing size of 2 should insert
	// from 99->98->etc.)
	mapInd, _ := strconv.Atoi(strings.Repeat("9", parsingSize))
	for VarMap[uint8(mapInd)] != nil {
		mapInd -= 1
	}

	// Finalize and insert the new var into the previously determined index
	VarMap[uint8(mapInd)] = &Y2KVar{
		ID:        uint8(mapInd),
		Size:      uint8(len(input)),
		StringVal: input,
		NumberVal: utils.SafeStrToInt(input),
		Type:      argType,
	}
}
