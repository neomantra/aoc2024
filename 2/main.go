// https://adventofcode.com/2024/day/2
// go run 2/main.go 2/2.txt

package main

import (
	"bufio"
	"fmt"
	"os"
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

type Report []int64

func isReportSafe(r Report) bool {
	var hasBeenIncreasing bool
	const minDiff, maxDiff = 1, 3
	for i := 1; i < len(r); i++ { // note starting at 1
		diff := r[i] - r[i-1]

		// check adjacent level threshold
		absDiff := abs(diff)
		if absDiff < minDiff || absDiff > maxDiff {
			return false
		}

		// we are not safe if we don't have same trend
		isIncreasingNow := (diff > 0)
		if i == 1 {
			// set initial condition
			hasBeenIncreasing = isIncreasingNow
		}
		if isIncreasingNow != hasBeenIncreasing {
			return false
		}
	}
	return true
}

func main() {
	// Open and read data file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	var reports []Report
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var report []int64
		strs := strings.Fields(scanner.Text())
		for _, str := range strs {
			v := stripError(strconv.ParseInt(str, 10, 64))
			report = append(report, v)
		}

		reports = append(reports, report)
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", scanner.Err().Error())
		os.Exit(1)
	}

	// part 1
	safeCount := 0
	for _, report := range reports {
		if isReportSafe(report) {
			safeCount++
		}
	}
	fmt.Println("2.1:", safeCount)
}
