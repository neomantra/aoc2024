// https://adventofcode.com/2024/day/17
// go run 17/main.go 17/17.txt

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

type Machine struct {
	A, B, C int // registers
	I       int // instruction pointer
	Program []int
	Output  []int
}

func NewMachine(puzzle string) *Machine {
	// split puzzle parts
	parts := strings.Split(puzzle, "\n\n")
	registers, program := parts[0], parts[1][:]

	var regRegexp = regexp.MustCompile(`\D*(\d*)\D*(\d*)\D*(\d*)`)
	match := regRegexp.FindStringSubmatch(registers)
	if len(match) != 4 {
		return nil
	}

	machine := Machine{}
	machine.A, _ = strconv.Atoi(match[1])
	machine.B, _ = strconv.Atoi(match[2])
	machine.C, _ = strconv.Atoi(match[3])
	machine.I = 0

	codes := strings.Split(program, ": ")
	if len(codes) != 2 {
		return nil
	}
	for _, c := range strings.Split(codes[1], ",") {
		n, _ := strconv.Atoi(c)
		machine.Program = append(machine.Program, n)
	}
	return &machine
}

func (m *Machine) View() string {
	var output, outputNum string
	for i, o := range m.Output {
		outputNum += string(o + 48)
		output += string(o + 48)
		if i != len(m.Output)-1 {
			output += ","
		}
	}
	return fmt.Sprintf("A=%d B=%d C=%d I=%d\nProgram: %v\nOutTxt: %s\n",
		m.A, m.B, m.C, m.I, m.Program, output)
}

func (m *Machine) OutputNumView() string {
	var outputNum string
	for _, o := range m.Output {
		outputNum += string(o + 48)
	}
	return outputNum
}

///////////////////////////////////////////////////////////////////////////////

// Returns false if halted (true means can keep running)
func (m *Machine) Step() bool {
	if m.I+1 >= len(m.Program) {
		return false
	}
	opCode, literal := m.Program[m.I], m.Program[m.I+1]

	var combo int
	switch literal {
	case 0:
		combo = literal
	case 1:
		combo = literal
	case 2:
		combo = literal
	case 3:
		combo = literal
	case 4:
		combo = m.A
	case 5:
		combo = m.B
	case 6:
		combo = m.C
	case 7:
		// no combo
	default:
		return false // invalid operation
	}

	incrI := 2
	switch opCode {
	case 0: // adv
		m.A = m.A / (1 << combo) // 2^combo
	case 1: // xor
		m.B = m.B ^ literal
	case 2: // bst
		m.B = combo % 8
	case 3: // jnz
		if m.A != 0 {
			m.I = literal
			incrI = 0
		}
	case 4: // bxc
		m.B = m.B ^ m.C
	case 5: // out
		m.Output = append(m.Output, combo%8)
	case 6: // bdv
		m.B = m.A / (1 << combo) // 2^combo
	case 7: // cdv
		m.C = m.A / (1 << combo) // 2^combo
	default:
		return false // invalid opcde
	}

	m.I += incrI
	return true
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	machineData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	machine := NewMachine(string(machineData))
	if machine == nil {
		fmt.Fprintf(os.Stderr, "Bad machine data\n")
		os.Exit(1)
	}
	fmt.Println(machine.View())
	for machine.Step() {
		// run till halt
		fmt.Println(machine.View())
	}
	fmt.Println("\n15.1: ", machine.View(), "\n", machine.OutputNumView())
}
