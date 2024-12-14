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

///////////////////////////////////////////////////////////////////////////////

// Finds the cheapest play... returns 0 if cannot win
func (g ClawGame) CheapestPlay() int {
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

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	clawGameData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	games := NewClawGames(string(clawGameData))
	if games == nil {
		fmt.Fprintf(os.Stderr, "Bad claw game data\n")
		os.Exit(1)
	}

	// part 1
	cost := 0
	for _, game := range games {
		thisCost := game.CheapestPlay()
		cost += thisCost
		fmt.Println(game.String(), thisCost)
	}
	fmt.Println("13.1:", cost)
}
