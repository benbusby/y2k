package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"y2k/src/utils"
)

type Y2KWhile struct {
	VarID      uint8
	CompSize   uint8
	CompValue  string
	Comparison func(*Y2KVar, []string) bool
}

// ParseWhile compares a variable against a raw value and parses a segment of
// the timestamp until the comparison is false. The segment of the timestamp
// used for the loop is determined by a function terminator ("1999") or the end
// of the timestamp if the terminator is not found.
func (y2k *Y2K) ParseWhile(timestamp string, y2kWhile Y2KWhile) string {
	command, _ := strconv.Atoi(timestamp[:y2k.Digits])

	if y2kWhile.VarID == 0 {
		y2kWhile.VarID = uint8(command)
		y2k.DebugMsg(4, fmt.Sprintf("Variable ID: %d", y2kWhile.VarID))
	} else if y2kWhile.Comparison == nil {
		y2kWhile.Comparison = ComparisonMap[uint8(command)]
		y2k.DebugMsg(4, fmt.Sprintf("Comp Function: %d", command))
	} else if y2kWhile.CompSize == 0 {
		y2kWhile.CompSize = uint8(command)
		y2k.DebugMsg(4, fmt.Sprintf("Comp Value Size: %d", command))
	} else {
		strVal := strconv.Itoa(command)
		y2kWhile.CompValue += strVal
		y2k.DebugMsg(6, fmt.Sprintf("(+ value: %s)", strVal))

		if len(y2kWhile.CompValue) >= int(y2kWhile.CompSize) {
			targetVar := VarMap[y2kWhile.VarID]

			// Comparison functions need the raw comparison value passed to
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
				y2kWhile.CompValue[:y2kWhile.CompSize],
				y2k.Digits)

			// Extract the index of the loop terminator and the subset of the
			// timestamp that should be returned back to the main interpreter
			// loop.
			timestampFnTerm := strings.Index(timestamp, utils.FnTerm)
			nextIterTimestamp := timestamp[timestampFnTerm+len(utils.FnTerm):]

			// If there isn't a function terminator, assume that the while loop
			// terminates at the end of the timestamp.
			if timestampFnTerm < 0 {
				timestampFnTerm = len(timestamp)
				nextIterTimestamp = ""
			}

			// Determine the segment of the timestamp that will be parsed on
			// each iteration of the while loop.
			whileTimestamp := timestamp[y2k.Digits:timestampFnTerm]
			y2k.DebugMsg(0, utils.DebugDivider)

			for y2kWhile.Comparison(targetVar, splitComp) {
				y2k.DebugMsg(0, fmt.Sprintf("RUN WHILE LOOP: %s\n%s",
					whileTimestamp,
					utils.DebugDivider))
				y2k.Parse(whileTimestamp)
				y2k.DebugMsg(0, utils.DebugDivider)
			}

			return nextIterTimestamp
		}
	}

	return y2k.ParseWhile(timestamp[y2k.Digits:], y2kWhile)
}
