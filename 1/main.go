// https://adventofcode.com/2024/day/1
// go run 1/main.go 1/1.txt

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func stripError[T any](result T, _ error) T {
	return result
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func count(x int64, s []int64) int64 {
	// could be more efficient since it is sorted but meh
	var count int64 = 0
	for _, v := range s {
		if v == x {
			count++
		}
	}
	return count
}

func main() {
	// Open and read data file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	var l, r []int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs := strings.Fields(scanner.Text())
		l = append(l, stripError(strconv.ParseInt(strs[0], 10, 64)))
		r = append(r, stripError(strconv.ParseInt(strs[1], 10, 64)))
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", scanner.Err().Error())
		os.Exit(1)
	}

	// part 1
	slices.Sort(l)
	slices.Sort(r)

	var totalDist int64 = 0
	for i := 0; i < len(l); i++ {
		totalDist += abs(r[i] - l[i])
	}
	fmt.Println("1.1:", totalDist)

	// part 2
	var similarityScore int64 = 0
	for _, v := range l {
		c := count(v, r)
		similarityScore += v * c
	}
	fmt.Println("1.2:", similarityScore)
}
