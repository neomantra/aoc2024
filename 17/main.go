// https://adventofcode.com/2024/day/17
// go run 17/main.go 17/17.txt

package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func WithCommas(nums []int) string {
	var str string
	for i, n := range nums {
		str += string(n + 48)
		if i != len(str)-1 {
			str += ","
		}
	}
	return str
}

///////////////////////////////////////////////////////////////////////////////

type Machine struct {
	A, B, C int // registers
	I       int // instruction pointer
	Program []int
	Output  []int
	StartA  int
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
	machine.StartA = machine.A

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

func (m *Machine) Clone() *Machine {
	machine := Machine{A: m.A, B: m.B, C: m.C, I: m.I}
	for _, c := range m.Program {
		machine.Program = append(machine.Program, c)
	}
	for _, c := range m.Output {
		machine.Output = append(machine.Output, c)
	}
	return &machine
}

func (m Machine) View() string {
	return fmt.Sprintf("A=%d B=%d C=%d I=%d\nProgram: %v\nOutTxt: %s\n",
		m.A, m.B, m.C, m.I, m.Program, WithCommas(m.Output))
}

///////////////////////////////////////////////////////////////////////////////

func OpcodeToString(opCode int) string {
	switch opCode {
	case 0:
		return "adv"
	case 1:
		return "bxl"
	case 2:
		return "bst"
	case 3:
		return "jnz"
	case 4:
		return "bxc"
	case 5:
		return "out"
	case 6:
		return "bdv"
	case 7:
		return "cdv"
	default:
		return "" // invalid opcde
	}
}

func OperandToString(operCode int) string {
	switch operCode {
	case 0:
		return strconv.Itoa(operCode)
	case 1:
		return strconv.Itoa(operCode)
	case 2:
		return strconv.Itoa(operCode)
	case 3:
		return strconv.Itoa(operCode)
	case 4:
		return "A"
	case 5:
		return "B"
	case 6:
		return "C"
	case 7:
		return "-"
	default:
		return "-" // invalid opcde
	}
}

func Dissassemble(opCode int, operCode int) string {
	if opCode == 1 || opCode == 3 || opCode == 4 {
		return OpcodeToString(opCode) + " " + strconv.Itoa(operCode)
	}
	return OpcodeToString(opCode) + " " + OperandToString(operCode)
}

func DetailDissassemble(opCode int, operCode int) string {
	var str string
	switch opCode {
	case 0: // adv
		str = fmt.Sprintf("A = A >> %s", OperandToString(operCode))
	case 1: // bxl
		str = fmt.Sprintf("B = B ^ %d", operCode)
	case 2: // bst
		str = fmt.Sprintf("B = %s %% 8", OperandToString(operCode))
	case 3: // jnz
		str = fmt.Sprintf("if A=!0 {jmp %d}", operCode)
	case 4: // bxc
		str = fmt.Sprintf("B = B ^ C")
	case 5: // out
		str = fmt.Sprintf("out %s %% 8", OperandToString(operCode))
	case 6: // bdv
		str = fmt.Sprintf("B = A >> %s", OperandToString(operCode))
	case 7: // cdv
		str = fmt.Sprintf("C = A >> %s", OperandToString(operCode))
	default:
		str = "--" // invalid opcde
	}
	return fmt.Sprintf("%-32s", str)
}

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
	case 1: // bxl
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

func (m *Machine) Run() {
	for m.Step() {
	}
}

///////////////////////////////////////////////////////////////////////////////

// search backwards for quine possiblities
func (m *Machine) QuineSearch() int {
	// we go backwards
	mods := []int{0, 1, 2, 3, 4, 5, 6, 7}
	alist := []int{0, 1, 2, 3, 4, 5, 6, 7}
	for i := len(m.Program) - 1; i >= 0; i-- {
		target := m.Program[i]
		var newAlist []int
		for _, a := range alist {
			for _, mod := range mods {
				newA := (a << 3) + mod
				mc := m.Clone()
				mc.A = newA
				mc.Run()
				if mc.Output[0] == target {
					newAlist = append(newAlist, newA)
				}
			}
		}
		alist = newAlist
	}

	// now we have alist of possible A values
	smallest := math.MaxInt
	for _, a := range alist {
		if a < smallest {
			smallest = a
		}
	}
	return smallest
}

///////////////////////////////////////////////////////////////////////////////

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Step key.Binding
	Go   key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Go, k.Step}
}

// FullHelp returns keybindings for the expanded help view. It's part of the key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Go, k.Step},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Go: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go"),
	),
	Step: key.NewBinding(
		key.WithKeys("s", " "),
		key.WithHelp("s", "step"),
	),
}

///////////////////////////////////////////////////////////////////////////////

type TModel struct {
	m    *Machine
	keys keyMap
	help help.Model
}

func NewTModel(m *Machine) *TModel {
	return &TModel{
		m:    m,
		keys: keys,
		help: help.New(),
	}
}

func (tm *TModel) Init() tea.Cmd {
	return nil
}

func (tm *TModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		tm.help.Width = msg.Width
		return tm, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tm.keys.Quit):
			return tm, tea.Quit
		case key.Matches(msg, tm.keys.Go):
			for tm.m.Step() {
			}
			return tm, tea.Quit
		case key.Matches(msg, tm.keys.Step):
			if !tm.m.Step() {
				// halted
				return tm, tea.Quit
			}
		}
	}
	return tm, nil
}

func (tm TModel) View() string {
	m := tm.m
	border := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
	inverted := lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000"))

	registerView := fmt.Sprintf("A | %16d | % 48s\nB | %16d | % 48s\nC | %16d | % 48s\nI | %16d |",
		m.A, strconv.FormatInt(int64(m.A), 2),
		m.B, strconv.FormatInt(int64(m.B), 2),
		m.C, strconv.FormatInt(int64(m.C), 2),
		m.I)

	var programView string
	for i := 0; i+1 < len(m.Program); i += 2 {
		op, oper := m.Program[i], m.Program[i+1]
		line := fmt.Sprintf("%d  %d | %s | %s",
			op, oper, Dissassemble(op, oper), DetailDissassemble(op, oper))
		if i == m.I {
			line = inverted.Render(line)
		}
		programView += line + "\n"
	}
	programView = programView[:len(programView)-1]

	return lipgloss.JoinVertical(lipgloss.Left,
		"Machine - "+strconv.Itoa(m.StartA),
		border.Render(registerView),
		"Program",
		border.Render(programView),
		"Output",
		border.Render(WithCommas(m.Output)),
		tm.help.View(tm.keys))
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

	// part 2
	machine = NewMachine(string(machineData))
	if len(os.Args) > 2 {
		a, err := strconv.Atoi(os.Args[2])
		if err == nil && a != 0 {
			machine.A = a
			machine.StartA = a
		}
	}

	aval := machine.QuineSearch()
	fmt.Print("\n15.2: ", aval, "\n\n")

	tm := NewTModel(machine)
	tm.m.A = aval
	if _, err := tea.NewProgram(tm).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
