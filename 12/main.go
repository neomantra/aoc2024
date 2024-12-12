// https://adventofcode.com/2024/day/12
// go run 10/main.go 12/12.txt

package main

import (
	"fmt"
	"os"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

type Point struct{ X, Y int }

type Garden struct {
	puzzle string
	rows   [][]byte // Y, X
	marks  [][]bool // Y, X
	extent int      // square
}

func NewGarden(puzzle string) *Garden {
	garden := Garden{puzzle: puzzle}

	// extract all the lines
	for _, row := range strings.Split(puzzle, "\n") {
		garden.rows = append(garden.rows, []byte(row))
	}
	garden.extent = len(garden.rows)
	garden.clearMarks()
	return &garden
}

func (g *Garden) IsInBounds(pt Point) bool {
	return pt.X >= 0 && pt.Y >= 0 && pt.X < g.extent && pt.Y < g.extent
}

func (g *Garden) GetPlant(pt Point) byte {
	if !g.IsInBounds(pt) {
		return 0
	}
	return g.rows[pt.Y][pt.X]
}

func (g *Garden) clearMarks() {
	g.marks = nil
	for _, row := range g.rows {
		g.marks = append(g.marks, make([]bool, len(row)))
	}

}

func (g *Garden) Mark(pt Point) {
	g.marks[pt.Y][pt.X] = true
}

func (g *Garden) IsMarked(pt Point) bool {
	if !g.IsInBounds(pt) {
		return true // OOB is always marked
	}
	return g.marks[pt.Y][pt.X]
}

func (g *Garden) View() string {
	var sb strings.Builder
	for _, row := range g.rows {
		sb.Write(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (g *Garden) IsSamePlant(otherPlant byte, pt Point) bool {
	thisPlant := g.GetPlant(pt)
	return thisPlant == otherPlant
}

func (g *Garden) IsBoundary(otherPlant byte, pt Point) bool {
	return !g.IsSamePlant(otherPlant, pt)
}

func (g *Garden) GetCellMetric(pt Point) RegionMetric {
	if !g.IsInBounds(pt) {
		return RegionMetric{}
	}
	metric := RegionMetric{Area: 1}
	thisPlant := g.rows[pt.Y][pt.X]
	if g.IsBoundary(thisPlant, Point{pt.X - 1, pt.Y + 0}) {
		metric.Perimeter++
	}
	if g.IsBoundary(thisPlant, Point{pt.X + 1, pt.Y + 0}) {
		metric.Perimeter++
	}
	if g.IsBoundary(thisPlant, Point{pt.X + 0, pt.Y - 1}) {
		metric.Perimeter++
	}
	if g.IsBoundary(thisPlant, Point{pt.X + 0, pt.Y + 1}) {
		metric.Perimeter++
	}
	return metric
}

///////////////////////////////////////////////////////////////////////////////

type RegionMetric struct {
	Area, Perimeter int
}

func (rm RegionMetric) Add(other RegionMetric) RegionMetric {
	return RegionMetric{
		Area:      rm.Area + other.Area,
		Perimeter: rm.Perimeter + other.Perimeter,
	}
}

func (rm RegionMetric) Cost() int {
	return rm.Area * rm.Perimeter
}

///////////////////////////////////////////////////////////////////////////////

func (g *Garden) regionScan(prevPlant byte, pt Point) RegionMetric {
	// if this point is not in bounds or different, return nil metric
	if g.IsMarked(pt) {
		return RegionMetric{}
	}
	plant := g.GetPlant(pt)
	if plant == 0 || prevPlant != plant {
		return RegionMetric{}
	}

	metric := g.GetCellMetric(pt)
	g.Mark(pt)
	metric = metric.Add(g.regionScan(prevPlant, Point{pt.X + 1, pt.Y + 0})) // right
	metric = metric.Add(g.regionScan(prevPlant, Point{pt.X - 1, pt.Y + 0})) // left
	metric = metric.Add(g.regionScan(prevPlant, Point{pt.X + 0, pt.Y - 1})) // up
	metric = metric.Add(g.regionScan(prevPlant, Point{pt.X + 0, pt.Y + 1})) // down
	return metric
}

func (g *Garden) TotalCost() int {
	var totalCost int
	g.clearMarks()
	for y := 0; y < g.extent; y++ {
		for x := 0; x < g.extent; x++ {
			pt := Point{x, y}
			if g.IsMarked(pt) {
				continue
			}

			// accumulate the metrics
			plant := g.GetPlant(pt)
			metric := g.regionScan(plant, pt)
			totalCost += metric.Cost()
		}
	}
	return totalCost
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	gardenData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	garden := NewGarden(string(gardenData))
	if garden == nil {
		fmt.Fprintf(os.Stderr, "Bad garden data\n")
		os.Exit(1)
	}

	// part 1
	fmt.Println(garden.View())
	totalCost := garden.TotalCost()
	fmt.Println("12.1:", totalCost)
}
