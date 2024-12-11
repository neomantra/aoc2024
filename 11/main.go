// https://adventofcode.com/2024/day/10
// go run 10/main.go 10/10.txt

package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func tenToPower(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////

type StoneRow struct {
	puzzle string
	stones []int
}

func NewStoneRow(puzzle string) *StoneRow {
	stoneRow := StoneRow{puzzle: puzzle}

	// extract all the lines
	for _, field := range strings.Fields(puzzle) {
		num, _ := strconv.Atoi(field)
		stoneRow.stones = append(stoneRow.stones, num)
	}
	return &stoneRow
}

func (sr *StoneRow) Blink() {
	lastIndex := len(sr.stones)
	for i := 0; i < lastIndex; i++ {
		stone := sr.stones[i]
		if stone == 0 {
			sr.stones[i] = 1
			continue
		}

		num, numDigits := stone, 0
		for num > 0 {
			num /= 10
			numDigits++
		}

		if numDigits%2 != 0 {
			sr.stones[i] *= 2024
		} else {
			// split odd stone in half
			factor := tenToPower(numDigits / 2)
			a := stone / factor
			b := stone - (a * factor)
			sr.stones[i] = a
			sr.stones = slices.Insert(sr.stones, i+1, b)
			lastIndex++
			i++
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

func (sr *StoneRow) View() string {
	var sb strings.Builder
	for _, stone := range sr.stones {
		sb.Write([]byte(strconv.Itoa(stone)))
		sb.WriteByte(' ')
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	stoneData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	stoneRow := NewStoneRow(string(stoneData))
	if stoneRow == nil {
		fmt.Fprintf(os.Stderr, "Bad stone data\n")
		os.Exit(1)
	}

	// part 1
	fmt.Println(stoneRow.View())
	for i := 0; i < 25; i++ {
		stoneRow.Blink()
		fmt.Println(stoneRow.View())
	}
	fmt.Println("11.1:", len(stoneRow.stones))
}
