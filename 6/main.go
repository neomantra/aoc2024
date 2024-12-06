// https://adventofcode.com/2024/day/6
// go run 6/main.go 6/6.txt

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

type Point struct{ X, Y int }

type Maze struct {
	Floorplan [][]byte // Y is array index, X is byte index sp [Y][X]
	Coloring  [][]byte
	Extent    Point
	GuardPos  Point
}

///////////////////////////////////////////////////////////////////////////////

const (
	GuardUp    = byte('^')
	GuardDown  = byte('v')
	GuardLeft  = byte('<')
	GuardRight = byte('>')
	Obstacle   = byte('#')
	Emptiness  = byte('.') // Buddha says that the emptiness of the maze is what gives it its function.
	Blank      = byte(' ')
	XSpot      = byte('X')
)

var guardRunes = []byte{GuardUp, GuardDown, GuardLeft, GuardRight}

func isGuard(c byte) bool {
	if c == GuardUp || c == GuardDown || c == GuardLeft || c == GuardRight {
		return true
	}
	return false
}

func isObstacle(c byte) bool {
	return c == Obstacle
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

func NewMaze(mazeStr string) *Maze {
	maze := &Maze{
		GuardPos: Point{X: -1, Y: -1},
	}

	for lineNum, line := range strings.Split(mazeStr, "\n") {
		maze.Extent.Y++
		maze.Extent.X = len(line) // assume all lines are same length
		maze.Floorplan = append(maze.Floorplan, []byte(line))
		maze.Coloring = append(maze.Coloring, bytes.Repeat([]byte{Emptiness}, len(line)))

		// find guard on floorplan
		if i := findGuard(line); i != -1 {
			maze.GuardPos = Point{X: i, Y: lineNum}
		}
	}

	if maze.GuardPos.X == -1 || maze.GuardPos.Y == -1 {
		return nil // no guard!
	}
	return maze
}

func (m *Maze) GetColorCount() int {
	count := 0
	for _, line := range m.Coloring {
		for _, c := range line {
			if c == XSpot {
				count++
			}
		}
	}
	return count
}

///////////////////////////////////////////////////////////////////////////////

// iterates the guard walking through the maze, coloring the map
func (m *Maze) WalkGuardAndColor() {
	// mark current position
	m.Coloring[m.GuardPos.Y][m.GuardPos.X] = 'X'
	done := false
	for !done {
		// look at the guard position
		gpos := m.GuardPos
		guardChar := m.Floorplan[gpos.Y][gpos.X]
		switch guardChar {
		case GuardUp:
			gpos.Y--
		case GuardDown:
			gpos.Y++
		case GuardLeft:
			gpos.X--
		case GuardRight:
			gpos.X++
		default:
			panic("bad guard position")
		}
		if !m.IsInBounds(gpos) {
			return // guard exited
		} else if m.IsObstacle(gpos) {
			// guard hit an obstacle, so turn right 90 and continue
			m.Floorplan[m.GuardPos.Y][m.GuardPos.X] = RotateGuard(guardChar)
			// TODO track infinite loops
		} else {
			// guard advances...
			// Empty current spot, color new spot, and move guard
			m.Floorplan[m.GuardPos.Y][m.GuardPos.X] = Emptiness
			m.Coloring[gpos.Y][gpos.X] = XSpot
			m.Floorplan[gpos.Y][gpos.X] = guardChar
			m.GuardPos = gpos
		}
	}
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
	for _, line := range maze.Coloring {
		fmt.Printf("%s\n", line)
	}
	fmt.Println("6.1:", maze.GetColorCount())
}
