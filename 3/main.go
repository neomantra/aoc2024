// https://adventofcode.com/2024/day/3
// go run 3/main.go 3/3.txt

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MulOp struct {
	A, B int
}

// retruns (a, b, error) with non-nil error on error
// currently not strict to digit spec
func extractMulArgs(mul string) (int, int, error) {
	strs := strings.Split(mul, ",")
	if len(strs) != 2 {
		return 0, 0, fmt.Errorf("invalid mul args: %s", mul)
	}
	a, err := strconv.Atoi(strs[0])
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.Atoi(strs[1])
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

// collectMulOpsextracts all mul(*) operators from source
func collectMulOps(src string) []MulOp {
	var result []MulOp
	var mulRegex = regexp.MustCompile(`mul\((.*?)\)`)
	matches := mulRegex.FindAllStringSubmatch(src, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		inParens := match[1]
		a, b, err := extractMulArgs(inParens)
		if err == nil {
			result = append(result, MulOp{A: a, B: b})
		} else {
			// try finding within the parens on failure
			// NOTE: need to add back closing )
			result = append(result, collectMulOps(inParens+")")...)
		}
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	sumResult := 0
	mulOps := collectMulOps(string(source))
	for _, mulOp := range mulOps {
		if err == nil {
			sumResult += (mulOp.A * mulOp.B)
		}
	}
	fmt.Println("3.1:", sumResult)

}
