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

///////////////////////////////////////////////////////////////////////////////

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

// collectMulOps extracts all valid mul(*) operators from source
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

// collectMulOpsDoDont extracts all valid mul(*) operators from source, with Do/Dont log
// Returns the ops and the next dontMultiply flag
func collectMulOpsDoDont(src string) []MulOp {
	var mulRegex = regexp.MustCompile(`^mul\((\d{1,3})\,(\d{1,3})\)`)
	var dontMultiply bool = false // default false means multiply-enabled
	var result []MulOp
	for i := 0; i < len(src); {
		// not the most efficient, but get 'er done
		doIdx := strings.Index(src[i:], "do()")
		if doIdx == 0 {
			dontMultiply = false
			i = i + len("do()")
			continue
		}

		dontIdx := strings.Index(src[i:], "don't()")
		if dontIdx == 0 {
			dontMultiply = true
			i = i + len("don't()")
			continue
		}

		mulIdx := strings.Index(src[i:], "mul(")
		if mulIdx == 0 {
			matches := mulRegex.FindStringSubmatch(src[i:])
			if len(matches) != 3 {
				i = i + len("mul(")
				continue
			}
			a, _ := strconv.Atoi(matches[1])
			b, _ := strconv.Atoi(matches[2])
			if dontMultiply == false {
				result = append(result, MulOp{A: a, B: b})
			}
			i = i + len(matches[0])
			continue
		}

		i = i + 1
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

	// part 2
	sumResult = 0
	mulOpsDoDont := collectMulOpsDoDont(string(source))
	for _, mulOp := range mulOpsDoDont {
		if err == nil {
			sumResult += (mulOp.A * mulOp.B)
		}
	}
	fmt.Println("3.2:", sumResult)
}
