// https://adventofcode.com/2024/day/15
// go run 15/main.go 15/15.txt

package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Empty = '.'
	Box   = 'O'
	Bot   = '@'
	Wall  = '#'
	LBox  = '['
	RBox  = ']'

	Up      = '^'
	Down    = 'v'
	Left    = '<'
	Right   = '>'
	Newline = '\n'
)

type Point struct{ X, Y int }

func (pt Point) Add(other Point) Point {
	return Point{pt.X + other.X, pt.Y + other.Y}
}

///////////////////////////////////////////////////////////////////////////////

type Warehouse struct {
	Map      [][]byte // Y, X
	Moves    []byte
	RobotPos Point
	Extent   Point
}

func NewWarehouse(puzzle string) *Warehouse {
	// split puzzle parts
	parts := strings.Split(puzzle, "\n\n")
	maze, moves := parts[0], parts[1][:]

	warehouse := Warehouse{}
	for y, line := range strings.Split(maze, "\n") {
		warehouse.Map = append(warehouse.Map, []byte(line))
		warehouse.Extent.X = len(line)
		for x, c := range line {
			if c == Bot {
				warehouse.RobotPos = Point{x, y}
			}
		}
	}
	warehouse.Extent.Y = len(warehouse.Map)
	warehouse.Moves = []byte(moves)
	return &warehouse
}

func (w *Warehouse) View() string {
	var sb strings.Builder
	for _, row := range w.Map {
		sb.Write(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////

func (w *Warehouse) IsInBounds(pt Point) bool {
	return pt.X >= 0 && pt.Y >= 0 && pt.X < w.Extent.X && pt.Y < w.Extent.Y
}

func (w *Warehouse) GetCell(pt Point) byte {
	if !w.IsInBounds(pt) {
		return Wall
	}
	return w.Map[pt.Y][pt.X]
}

func (w *Warehouse) SlideBox(boxPos Point, dir Point) bool {
	nextPos := boxPos.Add(dir)
	sideCell := w.GetCell(nextPos)
	if sideCell == Wall {
		return false
	} else if sideCell == Box {
		if !w.SlideBox(nextPos, dir) {
			return false
		}
	}

	w.Map[boxPos.Y][boxPos.X] = Empty
	w.Map[nextPos.Y][nextPos.X] = Box
	return true
}

// Moves the robot to the next position, according to rules
// Returns true if the robot can move, false otherwise
func (w *Warehouse) MoveRobot(dir Point) bool {
	// check next position
	nextPos := w.RobotPos.Add(dir)
	sideCell := w.GetCell(nextPos)
	if sideCell == Wall {
		return false
	} else if sideCell == Box {
		// check if the box can move
		if !w.SlideBox(nextPos, dir) {
			// could not move box, so don't move bot
			return false
		}
	} else if sideCell == LBox || sideCell == RBox {
		if !w.SlideExpandedBox(nextPos, dir) {
			// could not move box, so don't move bot
			return false
		}
	}

	w.Map[w.RobotPos.Y][w.RobotPos.X] = Empty
	w.Map[nextPos.Y][nextPos.X] = Bot
	w.RobotPos = nextPos
	return true
}

func (w *Warehouse) Operate() {
	// Runs all the robot moves
	for _, move := range w.Moves {
		switch move {
		case Newline:
			continue
		case Up:
			w.MoveRobot(Point{0, -1})
		case Down:
			w.MoveRobot(Point{0, +1})
		case Left:
			w.MoveRobot(Point{-1, 0})
		case Right:
			w.MoveRobot(Point{+1, 0})
		}
	}
}

// The GPS coordinate of a box is equal to 100 times its distance from
// the top edge of the map plus its distance from the left edge of the map.
func (w *Warehouse) GPSScore() int {
	score := 0
	for y := 0; y < w.Extent.Y; y++ {
		for x := 0; x < w.Extent.X; x++ {
			cell := w.Map[y][x]
			if cell == Box || cell == LBox {
				score += 100*y + x
			}
		}
	}
	return score
}

///////////////////////////////////////////////////////////////////////////////

func (w *Warehouse) Expand() {
	var newMap [][]byte // Y, X
	for y := 0; y < w.Extent.Y; y++ {
		var b []byte
		for x := 0; x < w.Extent.X; x++ {
			c := w.Map[y][x]
			c1, c2 := c, c
			if c == Box {
				c1, c2 = LBox, RBox
			} else if c == Bot {
				c1, c2 = Bot, Empty
				w.RobotPos = Point{x, y}
			}
			b = append(b, c1, c2)
		}
		newMap = append(newMap, b)
	}
	w.Map = newMap
	w.Extent.X *= 2
	w.RobotPos.X *= 2
}

func (w *Warehouse) TryMoveBoxVertically(boxPos Point, dir Point) bool {
	boxCell := w.GetCell(boxPos)
	if boxCell == Empty {
		return true
	} else if boxCell == Wall {
		return false
	}

	var otherBoxPos Point
	if boxCell == LBox {
		otherBoxPos = boxPos.Add(Point{+1, 0})
	} else { // boxCell == RBox
		otherBoxPos = boxPos.Add(Point{-1, 0})
	}

	nextPos, otherNextPos := boxPos.Add(dir), otherBoxPos.Add(dir)
	nextCell, otherNextCell := w.GetCell(nextPos), w.GetCell(otherNextPos)
	if nextCell == Wall || otherNextCell == Wall {
		return false
	}
	if (nextCell == Empty || w.TryMoveBoxVertically(nextPos, dir)) &&
		(otherNextCell == Empty || w.TryMoveBoxVertically(otherNextPos, dir)) {
		return true
	}

	return false
}

// we call this when we know the box can move
func (w *Warehouse) MoveBoxVertically(boxPos Point, dir Point) bool {
	boxCell := w.GetCell(boxPos)
	if boxCell == Empty {
		return true
	} else if boxCell == Wall {
		return false
	}

	var otherBoxPos Point
	var otherBoxCell byte
	if boxCell == LBox {
		otherBoxCell = RBox
		otherBoxPos = boxPos.Add(Point{+1, 0})
	} else { // boxCell == RBox
		otherBoxCell = LBox
		otherBoxPos = boxPos.Add(Point{-1, 0})
	}

	nextPos, otherNextPos := boxPos.Add(dir), otherBoxPos.Add(dir)
	nextCell, otherNextCell := w.GetCell(nextPos), w.GetCell(otherNextPos)
	if nextCell == Wall || otherNextCell == Wall {
		return false
	}
	if nextCell == LBox || nextCell == RBox {
		w.MoveBoxVertically(nextPos, dir)
	}
	if otherNextCell == LBox || otherNextCell == RBox {
		w.MoveBoxVertically(otherNextPos, dir)
	}
	w.Map[nextPos.Y][nextPos.X] = boxCell
	w.Map[otherNextPos.Y][otherNextPos.X] = otherBoxCell
	w.Map[boxPos.Y][boxPos.X] = Empty
	w.Map[otherBoxPos.Y][otherBoxPos.X] = Empty
	return false
}

func (w *Warehouse) SlideExpandedBox(boxPos Point, dir Point) bool {
	// moving left/right
	if dir.X != 0 {
		nextPos := boxPos.Add(dir)
		nextNextPos := nextPos.Add(dir)
		nextNextCell := w.GetCell(nextNextPos)
		if nextNextCell == Wall {
			return false
		}
		if (nextNextCell == LBox || nextNextCell == RBox) &&
			!w.SlideExpandedBox(nextNextPos, dir) {
			return false
		}

		thisCell := w.GetCell(boxPos)
		sideCell := w.GetCell(nextPos)
		w.Map[boxPos.Y][boxPos.X] = Empty
		w.Map[nextPos.Y][nextPos.X] = thisCell
		w.Map[nextNextPos.Y][nextNextPos.X] = sideCell
		return true
	}

	// moving up/down
	if w.TryMoveBoxVertically(boxPos, dir) {
		w.MoveBoxVertically(boxPos, dir)
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	warehouseData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	warehouse := NewWarehouse(string(warehouseData))
	if warehouse == nil {
		fmt.Fprintf(os.Stderr, "Bad warehouse data\n")
		os.Exit(1)
	}
	fmt.Print(warehouse.View(), "\n\n")
	warehouse.Operate()
	fmt.Print(warehouse.View(), "\n")
	fmt.Println("15.1:", warehouse.GPSScore())

	// part 2
	warehouse = NewWarehouse(string(warehouseData))
	if warehouse == nil {
		fmt.Fprintf(os.Stderr, "Bad warehouse data\n")
		os.Exit(1)
	}
	fmt.Print("\n\nPart 2\n", warehouse.View(), "\n\n")
	warehouse.Expand()
	fmt.Print(warehouse.View(), "\n")
	warehouse.Operate()
	fmt.Print(warehouse.View(), "\n")
	fmt.Println("15.2:", warehouse.GPSScore())
}
