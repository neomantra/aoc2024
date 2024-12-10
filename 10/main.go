// https://adventofcode.com/2024/day/10
// go run 10/main.go 10/10.txt

package main

import (
	"fmt"
	"os"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

type Point struct{ X, Y int }

type Island struct {
	puzzle  string
	topoMap [][]byte // (Y,X) height ('0'- '9')
	extent  Point
}

func NewIsland(puzzle string) *Island {
	isld := Island{puzzle: puzzle}

	// extract all the lines
	for _, row := range strings.Split(puzzle, "\n") {
		isld.topoMap = append(isld.topoMap, []byte(row))
	}
	if len(isld.topoMap) == 0 {
		return nil
	}

	isld.extent.X, isld.extent.Y = len(isld.topoMap[0]), len(isld.topoMap) // assumes uniform input

	return &isld
}

func (isld *Island) IsInBounds(p Point) bool {
	if p.X >= 0 && p.X < isld.extent.X && p.Y >= 0 && p.Y < isld.extent.Y {
		return true
	}
	return false
}

func (isld *Island) GetCellVal(pt Point) byte {
	return isld.topoMap[pt.Y][pt.X]
}

// Returns the total score from the hike point
// Returns 0 if the hike is unsuccessful
func (isld *Island) hikeStep(pt Point, prevVal byte) []Point {
	if !isld.IsInBounds(pt) {
		return nil
	}

	thisVal := isld.GetCellVal(pt)
	diff := thisVal - prevVal
	if diff != 1 {
		return nil // too big of a step or receding, so didin't make it
	}

	if thisVal == '9' {
		return []Point{pt} // we reached a height
	}

	var foundPoints []Point
	foundPoints = append(foundPoints, isld.hikeStep(Point{pt.X + 0, pt.Y - 1}, thisVal)...) // up
	foundPoints = append(foundPoints, isld.hikeStep(Point{pt.X + 0, pt.Y + 1}, thisVal)...) // down
	foundPoints = append(foundPoints, isld.hikeStep(Point{pt.X - 1, pt.Y + 0}, thisVal)...) // left
	foundPoints = append(foundPoints, isld.hikeStep(Point{pt.X + 1, pt.Y + 0}, thisVal)...) // right

	// simple unique
	foundPointMap := make(map[Point]bool)
	for _, pt := range foundPoints {
		foundPointMap[pt] = true
	}
	foundPoints = nil
	for k, _ := range foundPointMap {
		foundPoints = append(foundPoints, k)
	}
	return foundPoints
}

func (isld *Island) SumAllTrailheadScores() int {
	sum := 0
	for y := 0; y < len(isld.topoMap); y++ {
		row := isld.topoMap[y]
		for x := 0; x < len(row); x++ {
			cell := row[x]
			if cell == '0' {
				// only handle trailheads
				foundPoints := isld.hikeStep(Point{X: x, Y: y}, '0'-1)
				sum += len(foundPoints)
			}
		}
	}
	return sum
}

///////////////////////////////////////////////////////////////////////////////

func (isld *Island) TopoMapView() string {
	var sb strings.Builder
	for _, line := range isld.topoMap {
		sb.Write(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	islandData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	isld := NewIsland(string(islandData))
	if isld == nil {
		fmt.Fprintf(os.Stderr, "Bad island data\n")
		os.Exit(1)
	}

	// part 1
	fmt.Println(isld.TopoMapView())
	fmt.Println("10.1:", isld.SumAllTrailheadScores())
}
