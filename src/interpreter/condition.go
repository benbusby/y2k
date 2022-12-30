package interpreter

import (
	"fmt"
	"github.com/benbusby/y2k/src/utils"
	"reflect"
	"strings"
)

// ComparisonMap holds an int->function mapping to compare a variable against
// an arbitrary value.
var ComparisonMap = map[uint8]func(*Y2KVar, []string) bool{
	1: EqualTo,
	2: LessThan,
	3: GreaterThan,
	4: IsDivisible,
}

type Y2KCond struct {
	VarID       uint8
	CompFn      uint8
	Loop        bool
	CompValSize uint8
	value       string
}

// RunCond evaluates a timestamp as either a loop or a standalone "if" statement.
// All conditions compare a variable's value against a slice of strings, with the
// latter getting converted to the variable's data type in the comparison function.
// For example, if comparing an integer variable against ["8", "9"], the integer
// would need to have the number 89 stored as its numeric value. If comparing a
// string variable, it would need to have "hi" stored as its string value.
func (y2kComp Y2KCond) RunCond(
	y2k Y2K,
	timestamp string,
	target *Y2KVar,
	splitComp []string,
) bool {
	var result string

	if y2kComp.Loop {
		for ComparisonMap[y2kComp.CompFn](target, splitComp) {
			result = y2k.Parse(timestamp)
		}
	} else {
		if ComparisonMap[y2kComp.CompFn](target, splitComp) {
			result = y2k.Parse(timestamp)
		}
	}

	if result == utils.ContinueCmd {
		// Break out of timestamp if commanded
		return true
	}

	return false
}

// EqualTo checks string or numeric equality
func EqualTo(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return y2kVar.strVal == utils.StrArrToPrintable(values)
	case Y2KNumber:
		return y2kVar.intVal == utils.StrArrToInt(values)
	}

	return false
}

// LessThan checks if a string's length is less than the specified value,
// or if a number is less than a different numeric value.
func LessThan(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return len(y2kVar.strVal) < utils.StrArrToInt(values)
	case Y2KNumber:
		return y2kVar.intVal < utils.StrArrToInt(values)
	}

	return false
}

// GreaterThan checks if a string's length is greater than the specified value,
// or if a number is greater than a different numeric value.
func GreaterThan(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return len(y2kVar.strVal) > utils.StrArrToInt(values)
	case Y2KNumber:
		return y2kVar.intVal > utils.StrArrToInt(values)
	}

	return false
}

// IsDivisible checks if a numeric variable is evenly divisible by a
// specific number. Currently, there isn't an equivalent for string
// variables, so this will just return true in that case.
func IsDivisible(y2kVar *Y2KVar, values []string) bool {
	switch y2kVar.Type {
	case Y2KString:
		return true
	case Y2KNumber:
		return y2kVar.intVal%utils.StrArrToInt(values) == 0
	}

	return false
}

// ParseCondition compares a variable against a raw value and parses a segment of
// the timestamp until the comparison is false. The segment of the timestamp
// used for the loop is determined by a function terminator ("1999") or the end
// of the timestamp if the terminator is not found.
func (y2k Y2K) ParseCondition(timestamp string, val reflect.Value) string {
	y2kCond := val.Interface().(Y2KCond)

	input := timestamp[:y2k.Digits]
	y2k.DebugMsg(fmt.Sprintf("ParseCondition: [%s]%s",
		input,
		timestamp[y2k.Digits:]))

	y2kCond.value += input

	if len(y2kCond.value) >= int(y2kCond.CompValSize) {
		targetVar := GetVar(y2kCond.VarID)

		// CompFn functions need the raw comparison value passed to
		// them, because they treat values differently depending on the
		// target variable data type. It's easier to parse the comparison
		// value in as a string and then convert it back to a slice of
		// N-size strings than it is to create the slice during parsing, due
		// to differences in y2k.Digits values. For example -- parsing a 3
		// digit number "100XX..." with a 2-digit window would create a
		// slice of ["10", "0X"], where X is an unrelated digit for a
		// subsequent command. Parsing it as a string and then splitting it,
		// however, creates ["10", "0"].
		splitComp := utils.SplitStrByN(
			y2kCond.value[:y2kCond.CompValSize],
			y2k.Digits)

		// Extract the index of the cond terminator and the subset of the
		// timestamp that should be returned to the main interpreter loop.
		condTerm := utils.GetCondTerm(y2kCond.Loop)
		timestampFnTerm := strings.Index(timestamp, condTerm)
		nextIterTimestamp := timestamp[timestampFnTerm+len(condTerm)-1:]

		// If there isn't a function terminator, assume that the condition
		// terminates at the end of the timestamp.
		if timestampFnTerm < 0 {
			timestampFnTerm = len(timestamp)
			nextIterTimestamp = ""
		}

		// Determine the segment of the timestamp that will be parsed on
		// each iteration of the while loop.
		whileTimestamp := timestamp[y2k.Digits:timestampFnTerm]
		y2k.DebugMsg(utils.DebugDivider)

		stop := y2kCond.RunCond(y2k, whileTimestamp, targetVar, splitComp)

		// Conditions can optionally break out of the timestamp using the
		// CONTINUE command (see interpreter.go). In this instance, the next
		// timestamp passed back to the parser should be empty.
		if stop {
			return ""
		}

		return nextIterTimestamp
	}

	return y2k.ParseCondition(timestamp[y2k.Digits:], reflect.ValueOf(y2kCond))
}
