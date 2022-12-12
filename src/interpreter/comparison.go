package interpreter

import (
	"y2k/src/utils"
)

// ComparisonMap holds an int->function mapping to compare a variable against
// an arbitrary value.
var ComparisonMap = map[uint8]func(*Y2KVar, []string) bool{
	1: EqualTo,
	2: LessThan,
	3: GreaterThan,
}

// EqualTo checks string or numeric equality
func EqualTo(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return y2kVar.StringVal == utils.StrArrToPrintable(values)
	case Y2KNumber:
		return y2kVar.NumberVal == utils.SafeStrArrToInt(values)
	}

	return false
}

// LessThan checks if a string's length is less than the specified value,
// or if a number is less than a different numeric value.
func LessThan(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return len(y2kVar.StringVal) < utils.SafeStrArrToInt(values)
	case Y2KNumber:
		return y2kVar.NumberVal < utils.SafeStrArrToInt(values)
	}

	return false
}

// GreaterThan checks if a string's length is greater than the specified value,
// or if a number is greater than a different numeric value.
func GreaterThan(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return len(y2kVar.StringVal) > utils.SafeStrArrToInt(values)
	case Y2KNumber:
		return y2kVar.NumberVal > utils.SafeStrArrToInt(values)
	}

	return false
}
