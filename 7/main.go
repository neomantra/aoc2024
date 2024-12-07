// https://adventofcode.com/2024/day/7
// go run 7/main.go 7/7.txt

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// Operators are always evaluated left-to-right, not according to precedence rules.

type Equation struct {
	Result int
	Args   []int
}

func (e Equation) Calc(ops []Op) (int, error) {
	if len(ops) != len(e.Args)-1 {
		return 0, fmt.Errorf("Wrong number of operators")
	}
	result := e.Args[0] // seed with first arg
	for i := 0; i < len(e.Args)-1; i++ {
		result = ops[i].Apply(result, e.Args[i+1])
	}
	return result, nil
}

///////////////////////////////////////////////////////////////////////////////

type Op interface {
	Glyph() byte
	Apply(a, b int) int
}

///////////////////////////////////////////////////////////////////////////////

type AddOp struct{}

func (op AddOp) Glyph() byte { return '+' }

func (op AddOp) Apply(a, b int) int { return a + b }

///////////////////////////////////////////////////////////////////////////////

type MulOp struct{}

func (op MulOp) Glyph() byte { return '*' }

func (op MulOp) Apply(a, b int) int { return a * b }

///////////////////////////////////////////////////////////////////////////////

var allOps = []Op{AddOp{}, MulOp{}}
var allOpsBase = len(allOps)

// permuteOps returns a slice of operators that correspond to the given seed
// The seed is a base_len(allOps) number, where each digit represents an operator
func permuteOps(seed int, length int) []Op {
	if seed < 0 || length <= 0 {
		return nil
	}
	ceiling := int(math.Pow(float64(allOpsBase), float64(length)))
	if seed > ceiling {
		return nil // overflow
	}
	ops := make([]Op, length)
	for i := 0; i < length; i++ {
		allOpsIndex := seed % allOpsBase
		seed /= allOpsBase
		ops[i] = allOps[allOpsIndex]
	}
	return ops
}

///////////////////////////////////////////////////////////////////////////////

// Returns the set of operators that satisfy the Equation
// Returns nil if none exist
func FindOps(e Equation) []Op {
	if len(e.Args) == 0 {
		return nil
	}
	// we will permute all the operators, returning the first successful set
	i := 0
	for {
		ops := permuteOps(i, len(e.Args)-1)
		if ops == nil {
			return nil // we are out of permutations
		}
		result, err := e.Calc(ops)
		if err != nil {
			return nil
		}
		if result == e.Result {
			return ops
		}
		i += 1
	}
}

///////////////////////////////////////////////////////////////////////////////

// NewEquations parses a string of equations into a slice of Equation structs
// Returns nil if any line is malformed
func NewEquations(data string) []Equation {
	var equations []Equation
	for _, line := range strings.Split(data, "\n") {
		var equation Equation
		pair := strings.Split(line, ":")
		if len(pair) != 2 {
			return nil // bad line
		}
		equation.Result, _ = strconv.Atoi(pair[0])
		for _, field := range strings.Fields(pair[1]) {
			arg, _ := strconv.Atoi(field)
			equation.Args = append(equation.Args, arg)
		}
		equations = append(equations, equation)
	}
	return equations
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	equationData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	equations := NewEquations(string(equationData))
	if equations == nil {
		fmt.Fprintf(os.Stderr, "Bad equation data\n")
		os.Exit(1)
	}
	sum := 0
	for _, e := range equations {
		if ops := FindOps(e); ops != nil {
			sum += e.Result
		}
	}
	fmt.Println("7.1:", sum)
}
