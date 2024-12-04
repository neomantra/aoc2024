// https://adventofcode.com/2024/day/4
// go run 4/main.go 4/4.txt

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/NimbleMarkets/ntcharts/canvas"
)

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func maxOf(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func signOf(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

///////////////////////////////////////////////////////////////////////////////

type Board struct {
	puzzle string
	lines  []string
	counts [][]int // (y,x) matrix of hit counts (for coloring)

	canvas     canvas.Model // for plotting
	maxX, maxY int
}

func (b *Board) Canvas() *canvas.Model {
	return &b.canvas
}

func NewBoard(puzzle string) *Board {
	// extract all the lines
	lines := strings.Split(puzzle, "\n")
	if len(lines) == 0 {
		return nil
	}

	lenX, lenY := len(lines[0]), len(lines) // assumes uniform input
	canvas := canvas.New(lenX, lenY)
	canvas.SetLines(lines)

	// fill out counts matrix
	var counts [][]int
	for i := 0; i < lenY; i++ {
		counts = append(counts, make([]int, lenX)) // 0-filled
	}

	return &Board{
		puzzle: puzzle,
		lines:  lines,
		counts: counts,
		maxX:   maxOf(0, lenX-1),
		maxY:   maxOf(0, lenY-1),
		canvas: canvas,
	}
}

func (b *Board) CharAt(x int, y int) byte {
	if x < 0 || x > b.maxX || y < 0 || y > b.maxY {
		// out of bounds, return empty rune
		return 0
	}
	return b.lines[y][x]
}

func (b *Board) StringLine(x, y, length, xdir, ydir int) string {
	// build the string via iteration
	signX, signY := signOf(xdir), signOf(ydir)
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		if c := b.CharAt(x, y); c == 0 {
			// out of bounds, stop
			break
		} else {
			buffer.WriteByte(c)
		}
		x += signX
		y += signY
	}
	return buffer.String()
}

func (b *Board) CountAt(word string, x int, y int) int {
	if word == "" {
		return 0
	}
	if b.lines[y][x] != word[0] {
		return 0 // quick exit
	}

	var tests []string
	lenW := len(word)
	tests = append(tests, b.StringLine(x, y, lenW, -1, 0))  // left
	tests = append(tests, b.StringLine(x, y, lenW, +1, 0))  // right
	tests = append(tests, b.StringLine(x, y, lenW, 0, -1))  // up
	tests = append(tests, b.StringLine(x, y, lenW, 0, +1))  // down
	tests = append(tests, b.StringLine(x, y, lenW, -1, -1)) // up-left
	tests = append(tests, b.StringLine(x, y, lenW, +1, -1)) // up-right
	tests = append(tests, b.StringLine(x, y, lenW, -1, +1)) // down-left
	tests = append(tests, b.StringLine(x, y, lenW, +1, +1)) // down-right

	count := 0
	for _, t := range tests {
		if t == word {
			count += 1
		}
	}
	return count
}

func (b *Board) CountWord(word string) int {
	// we are going to find Xs and search from there.
	sum := 0
	for y := 0; y < b.maxY+1; y++ {
		for x := 0; x < b.maxX+1; x++ {
			num := b.CountAt(word, x, y)
			sum += num
		}
	}
	return sum
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	puzzle, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	board := NewBoard(string(puzzle))
	if board == nil {
		fmt.Fprintf(os.Stderr, "Error creating board\n")
		os.Exit(1)
	}
	// fmt.Println(board.Canvas().View())
	fmt.Println("3.2:", board.CountWord("XMAS"))
}
