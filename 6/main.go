// https://adventofcode.com/2024/day/6
// go run 6/main.go 6/6.txt

package main

import (
	"fmt"
	"os"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

type Point struct{ X, Y int }

type Maze struct {
	Floorplan [][]byte // Y is array index, X is byte index sp [Y][X]
	Coloring  [][]Color
	Extent    Point
	GuardPos  Point
}

///////////////////////////////////////////////////////////////////////////////

const (
	GuardUp     = byte('^')
	GuardDown   = byte('v')
	GuardLeft   = byte('<')
	GuardRight  = byte('>')
	Obstacle    = byte('#')
	Obstruction = byte('O')
	Emptiness   = byte('.') // Buddha says that the emptiness of the maze is what gives it its function.
	Blank       = byte(' ')
	HMove       = byte('-')
	VMove       = byte('|')
	HVMove      = byte('+')
)

var guardRunes = []byte{GuardUp, GuardDown, GuardLeft, GuardRight}

type Color byte

const (
	ColorNone  = 0
	ColorUp    = 1
	ColorDown  = 2
	ColorLeft  = 4
	ColorRight = 8
	ColorHoriz = ColorLeft | ColorRight
	ColorVert  = ColorUp | ColorDown
	ColorAll   = ColorUp | ColorDown | ColorLeft | ColorRight
)

func ColorFromGuard(g byte) Color {
	switch g {
	case GuardUp:
		return ColorUp
	case GuardDown:
		return ColorDown
	case GuardLeft:
		return ColorLeft
	case GuardRight:
		return ColorRight
	default:
		return ColorNone
	}
}

func (c Color) AsColorGlyph() byte {
	h := (c & ColorHoriz) != 0
	v := (c & ColorVert) != 0
	if h && !v {
		return HMove
	} else if !h && v {
		return VMove
	} else if h && v {
		return HVMove
	} else {
		return Blank
	}
}

func isGuard(c byte) bool {
	if c == GuardUp || c == GuardDown || c == GuardLeft || c == GuardRight {
		return true
	}
	return false
}

func isObstacle(c byte) bool {
	return c == Obstacle || c == Obstruction
}

// rotates the guard right (clockwise) 90 degrees
func RotateGuard(g byte) byte {
	switch g {
	case GuardUp:
		return GuardRight
	case GuardRight:
		return GuardDown
	case GuardDown:
		return GuardLeft
	case GuardLeft:
		return GuardUp
	default: // pass through bad char?
		return g
	}
}

///////////////////////////////////////////////////////////////////////////////

func (m *Maze) IsInBounds(p Point) bool {
	if p.X >= 0 && p.X < m.Extent.X && p.Y >= 0 && p.Y < m.Extent.Y {
		return true
	}
	return false
}

func (m *Maze) IsObstacle(p Point) bool {
	if p.X >= 0 && p.X < m.Extent.X && p.Y >= 0 && p.Y < m.Extent.Y {
		return isObstacle(m.Floorplan[p.Y][p.X])
	}
	return true // out of bounds is an obstacle
}

// Returns location of guard in line, or -1 if not found
func findGuard(line string) int {
	for i, c := range line {
		if isGuard(byte(c)) {
			return i
		}
	}
	return -1 // not found
}

///////////////////////////////////////////////////////////////////////////////

func (m *Maze) ClearColoring() {
	m.Coloring = nil
	for _, line := range m.Floorplan {
		m.Coloring = append(m.Coloring, make([]Color, len(line)))
	}
}

func NewMaze(mazeStr string) *Maze {
	maze := &Maze{
		GuardPos: Point{X: -1, Y: -1},
	}

	for lineNum, line := range strings.Split(mazeStr, "\n") {
		maze.Extent.Y++
		maze.Extent.X = len(line) // assume all lines are same length
		maze.Floorplan = append(maze.Floorplan, []byte(line))

		// find guard on floorplan
		if i := findGuard(line); i != -1 {
			maze.GuardPos = Point{X: i, Y: lineNum}
		}
	}

	if maze.GuardPos.X == -1 || maze.GuardPos.Y == -1 {
		return nil // no guard!
	}

	maze.ClearColoring()

	return maze
}

func (m *Maze) GetColorCount() int {
	count := 0
	for _, line := range m.Coloring {
		for _, c := range line {
			if c != ColorNone {
				count++
			}
		}
	}
	return count
}

// colors the maze at a point,
// assumes point is in-bounds
func (m *Maze) color(p Point, guardChar byte) Color {
	c := m.Coloring[p.Y][p.X]
	c |= ColorFromGuard(guardChar)
	m.Coloring[p.Y][p.X] = c
	return c
}

///////////////////////////////////////////////////////////////////////////////

// iterates the guard walking through the maze, coloring the map
// Returns false if there is a loop, true if the guard exits
func (m *Maze) WalkGuardAndColor() bool {
	// mark current position
	m.ClearColoring()
	spins := 0
	for spins < 4 {
		// look at the guard position
		guardChar := m.Floorplan[m.GuardPos.Y][m.GuardPos.X]
		newPos := m.GuardPos
		switch guardChar {
		case GuardUp:
			newPos.Y--
		case GuardDown:
			newPos.Y++
		case GuardLeft:
			newPos.X--
		case GuardRight:
			newPos.X++
		default:
			panic("bad guard position")
		}

		// bounds check
		if !m.IsInBounds(newPos) {
			return true // guard exited
		}

		// have we visited here before?
		guardColor := ColorFromGuard(guardChar)
		prevNewColor := m.Coloring[newPos.Y][newPos.X]
		if (guardColor | prevNewColor) == 0 {
			return false // we've been here before, so we're done
		}
		m.color(m.GuardPos, guardChar) // mark this square as visited

		// hit obstacle?
		if m.IsObstacle(newPos) {
			// guard hit an obstacle, so turn right 90 and continue
			m.Floorplan[m.GuardPos.Y][m.GuardPos.X] = RotateGuard(guardChar)
			spins += 1
		} else {
			// Guard advances...
			// Empty current spot, color new spot, and move guard
			m.Floorplan[m.GuardPos.Y][m.GuardPos.X] = Emptiness
			m.Floorplan[newPos.Y][newPos.X] = guardChar
			m.color(newPos, guardChar) // mark new square as visited
			m.GuardPos = newPos
			spins = 0
		}
	}
	return false // spinning endlessly
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	mazeData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	maze := NewMaze(string(mazeData))
	if maze == nil {
		fmt.Fprintf(os.Stderr, "Bad maze board\n")
		os.Exit(1)
	}

	// part 1
	maze.WalkGuardAndColor()
	for _, line := range maze.Floorplan {
		fmt.Printf("%s\n", line)
	}
	fmt.Println("")
	for _, colors := range maze.Coloring {
		for _, c := range colors {
			fmt.Printf("%c", c.AsColorGlyph())
		}
		fmt.Println()
	}
	fmt.Println("6.1:", maze.GetColorCount())
}
