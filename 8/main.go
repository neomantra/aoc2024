// https://adventofcode.com/2024/day/8
// go run 8/main.go 8/8.txt

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func maxOf(x, y int) int {
	if x > y {
		return x
	}
	return y
}

///////////////////////////////////////////////////////////////////////////////

const (
	EmptyGlyph    = '.'
	AntinodeGlyph = '#'
)

type Point struct{ X, Y int }

type City struct {
	puzzle     string
	antennas   [][]byte // Y, X
	antinodes  [][]byte // Y, X
	maxX, maxY int
}

func NewCity(puzzle string) *City {
	// extract all the lines
	antennas := bytes.Split([]byte(puzzle), []byte{'\n'})
	if len(antennas) == 0 {
		return nil
	}

	lenX, lenY := len(antennas[0]), len(antennas) // assumes uniform input

	c := &City{
		puzzle:   puzzle,
		antennas: antennas,
		maxX:     maxOf(0, lenX-1),
		maxY:     maxOf(0, lenY-1),
	}
	c.ClearAntinodes()
	return c
}

func (c *City) ClearAntinodes() {
	c.antinodes = nil
	for _, line := range c.antennas {
		emptyLine := strings.Repeat(string(EmptyGlyph), len(line))
		c.antinodes = append(c.antinodes, []byte(emptyLine))
	}
}

///////////////////////////////////////////////////////////////////////////////

func (c *City) setAntinode(x, y int) bool {
	if x < 0 || x > c.maxX || y < 0 || y > c.maxY {
		return false
	}
	c.antinodes[y][x] = AntinodeGlyph
	return true
}

func (c *City) MarkAntinodes(a, b Point, oneStep bool) {
	// an antinode occurs at any point that is perfectly in line with
	// two antennas of the same frequency - but only
	// when one of the antennas is twice as far away as the other.
	if !oneStep {
		// with harmonics enabled, start with antennae themselves
		c.setAntinode(a.X, a.Y)
		c.setAntinode(b.X, b.Y)

	}

	xDist, yDist := b.X-a.X, b.Y-a.Y
	ax, ay, bx, by := a.X, a.Y, b.X, b.Y
	for {
		ax -= xDist
		ay -= yDist
		bx += xDist
		by += yDist
		aInBounds := c.setAntinode(ax, ay)
		bInBounds := c.setAntinode(bx, by)
		if oneStep || !(aInBounds || bInBounds) {
			return
		}
	}

}

func (c *City) FindAntinodes(oneStep bool) {
	// march through each cell, if it's an antenna, try to find antinodes of all remaining antennas
	// we keep track of scanned frequencies to not double count
	c.ClearAntinodes()
	for y1 := 0; y1 <= c.maxY; y1++ {
		for x1 := 0; x1 <= c.maxX; x1++ {
			freq := c.antennas[y1][x1]
			if freq == EmptyGlyph {
				continue
			}

			// handle the rest of this row
			for x2 := x1 + 1; x2 <= c.maxX; x2++ {
				if c.antennas[y1][x2] == freq {
					c.MarkAntinodes(Point{x1, y1}, Point{x2, y1}, oneStep)
				}
			}
			// handle the rest of the rows
			for y2 := y1 + 1; y2 <= c.maxY; y2++ {
				for x2 := 0; x2 <= c.maxX; x2++ {
					if c.antennas[y2][x2] == freq {
						c.MarkAntinodes(Point{x1, y1}, Point{x2, y2}, oneStep)
					}
				}
			}
		}
	}
}

func (c *City) GetAntinodeCount() int {
	count := 0
	for _, line := range c.antinodes {
		for _, c := range line {
			if c != 0 && c != EmptyGlyph {
				count++
			}
		}
	}
	return count
}

func (c *City) View() string {
	var style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
	var antenna, antinodes strings.Builder
	for _, line := range c.antennas {
		antenna.Write(line)
		antenna.WriteByte('\n')
	}
	for _, line := range c.antinodes {
		antinodes.Write(line)
		antinodes.WriteByte('\n')
	}
	return lipgloss.JoinHorizontal(lipgloss.Left,
		style.Render(antenna.String()), style.Render(antinodes.String()))
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	cityData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	city := NewCity(string(cityData))
	if city == nil {
		fmt.Fprintf(os.Stderr, "Bad equation data\n")
		os.Exit(1)
	}

	// part 1
	city.FindAntinodes(true)
	fmt.Println(city.View())
	fmt.Println("8.1:", city.GetAntinodeCount())

	// part 2
	city = NewCity(string(cityData))
	city.FindAntinodes(false)
	fmt.Println(city.View())
	fmt.Println("8.2:", city.GetAntinodeCount())
}
