// https://adventofcode.com/2024/day/13
// go run 13/main.go 13/13.txt

package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	ButtonACost = 3
	ButtonBCost = 1
)

type Point struct{ X, Y int }

type ClawGame struct {
	ButtonA Point
	ButtonB Point
	Prize   Point
}

func readPoint(line string, re *regexp.Regexp) Point {
	match := re.FindStringSubmatch(line)
	x, _ := strconv.Atoi(match[1])
	y, _ := strconv.Atoi(match[2])
	return Point{x, y}
}

func NewClawGames(puzzle string) []ClawGame {
	var buttonARegexp = regexp.MustCompile(`Button A: X\+(\d*), Y\+(\d*)`)
	var buttonBRegexp = regexp.MustCompile(`Button B: X\+(\d*), Y\+(\d*)`)
	var prizeRegexp = regexp.MustCompile(`Prize: X=(\d*), Y=(\d*)`)
	games := []ClawGame{}
	lines := strings.Split(puzzle, "\n")
	for i := 0; i < len(lines); i += 4 {
		buttonA := readPoint(lines[i], buttonARegexp)
		buttonB := readPoint(lines[i+1], buttonBRegexp)
		prize := readPoint(lines[i+2], prizeRegexp)
		games = append(games, ClawGame{buttonA, buttonB, prize})
	}
	return games
}

func (g ClawGame) String() string {
	return fmt.Sprintf("{ A: %v, B: %v, P: %v", g.ButtonA, g.ButtonB, g.Prize)
}

func (g *ClawGame) ApplyConversion() {
	g.Prize.X += 10000000000000
	g.Prize.Y += 10000000000000
}

///////////////////////////////////////////////////////////////////////////////

// Finds the cheapest play... returns 0 if cannot win
func (g ClawGame) CheapestPlayBrute() int {
	// brute force since small range
	const MaxIters = 101
	minScore := math.MaxInt
	for a := 0; a < MaxIters; a++ {
		for b := 0; b < MaxIters; b++ {
			// does this button combo win?
			clawPt := Point{
				g.ButtonA.X*a + g.ButtonB.X*b,
				g.ButtonA.Y*a + g.ButtonB.Y*b}
			if clawPt == g.Prize {
				// yep, calculate cost and check if min
				thisCost := ButtonACost*a + ButtonBCost*b
				if thisCost < minScore {
					minScore = thisCost
				}
			}
		}
	}
	if minScore == math.MaxInt {
		return 0
	}
	return minScore
}

func (g ClawGame) CheapestPlayLinear() int {
	// matrix:
	//  |Ax Bx| |Ta| = |Px|
	//  |Ay By| |Tb| = |Py|
	Ax, Ay, Bx, By := g.ButtonA.X, g.ButtonA.Y, g.ButtonB.X, g.ButtonB.Y
	det := Ax*By - Ay*Bx
	if det == 0 {
		return 0 // no solution or we can't go backwards
	}
	// Cramer's rule
	Px, Py := g.Prize.X, g.Prize.Y
	Ta := (Px*By - Py*Bx) / det
	Tb := (Py*Ax - Px*Ay) / det

	// check if solution is valid
	x := Ta*Ax + Tb*Bx
	y := Ta*Ay + Tb*By
	if x == Px && y == Py {
		return ButtonACost*Ta + ButtonBCost*Tb
	} else {
		return 0
	}
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	clawGameData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	games := NewClawGames(string(clawGameData))
	if games == nil {
		fmt.Fprintf(os.Stderr, "Bad claw game data\n")
		os.Exit(1)
	}
	cost := 0
	for _, game := range games {
		thisCost := game.CheapestPlayBrute()
		cost += thisCost
	}
	fmt.Println("13.1:", cost)

	// part 2
	for i := 0; i < len(games); i++ {
		games[i].ApplyConversion()
	}
	cost = 0
	for _, game := range games {
		thisCost := game.CheapestPlayLinear()
		cost += thisCost
	}
	fmt.Println("13.2:", cost)
}
