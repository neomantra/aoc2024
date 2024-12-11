// https://adventofcode.com/2024/day/10
// go run 10/main.go 10/10.txt

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func tenToPower(n int) int {
	switch n {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	case 4:
		return 10000
	case 5:
		return 100000
	case 6:
		return 1000000
	}
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func countDigits(num int) int {
	numDigits := 0
	for num > 0 {
		num /= 10
		numDigits++
	}
	return numDigits
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
	newStones := make([]int, 0, len(sr.stones)*2)
	for i := 0; i < len(sr.stones); i++ {
		stone := sr.stones[i]
		if stone == 0 {
			newStones = append(newStones, 1)
			continue
		}

		numDigits := countDigits(stone)

		// if numDigits%2 != 0 {
		if (numDigits & 1) != 0 {
			newStones = append(newStones, stone*2024)
		} else {
			// split odd stone in half
			factor := tenToPower(numDigits / 2)
			a := stone / factor
			b := stone - (a * factor)
			newStones = append(newStones, a, b)
		}
	}
	sr.stones = newStones
}

///////////////////////////////////////////////////////////////////////////////

type Pair struct {
	Stone     int
	NumBlinks int
}

func NewPair(stone, numBlinks int) Pair { return Pair{Stone: stone, NumBlinks: numBlinks} }

var breakTimesMemo = make(map[Pair]int)

func breakTimes(stone int, numBlinks int) int {
	if numBlinks <= 0 {
		return 1
	}
	cachedCount, ok := breakTimesMemo[NewPair(stone, numBlinks)]
	if ok {
		return cachedCount
	}

	if stone == 0 {
		count := breakTimes(1, numBlinks-1)
		breakTimesMemo[NewPair(stone, numBlinks)] = count
		return count
	}

	numDigits := countDigits(stone)
	if (numDigits & 1) != 0 {
		count := breakTimes(stone*2024, numBlinks-1)
		breakTimesMemo[NewPair(stone, numBlinks)] = count
		return count
	} else {
		// split odd stone in half
		factor := tenToPower(numDigits / 2)
		a := stone / factor
		b := stone - (a * factor)
		count := breakTimes(a, numBlinks-1) + breakTimes(b, numBlinks-1)
		breakTimesMemo[NewPair(stone, numBlinks)] = count
		return count
	}
}

func (sr *StoneRow) CountAfterBlinking(numBlinks int) int {
	// go over each stone, processing it numBlinks times
	numStones := 0
	for i := 0; i < len(sr.stones); i++ {
		stone := sr.stones[i]
		numStones += breakTimes(stone, numBlinks)
	}
	return numStones
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
	count := stoneRow.CountAfterBlinking(25)
	fmt.Println("11.1:", count)

	// part 2
	stoneRow = NewStoneRow(string(stoneData))
	count = stoneRow.CountAfterBlinking(75)
	fmt.Println("11.2:", count)
}
