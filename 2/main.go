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

// isReportSafe returns true if the report is safe, false otherwise.
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

func dampenReport(r Report, pos int) Report {
	// remove element at pos
	// also return append(r[:pos], r[pos+1:]...)
	// but that looks weird in debugging
	var d Report
	for i, v := range r {
		if i != pos {
			d = append(d, v)
		}
	}
	return d
}

func isReportSafeDampened(r Report) bool {
	// if it is safe, then report so
	if isReportSafe(r) {
		return true
	}

	// we are allowed to be safe with a level removed
	// try every position, and we're safe if its safe
	for i := 0; i < len(r); i++ {
		dampenedReport := dampenReport(r, i)
		if isReportSafe(dampenedReport) {
			return true
		}
	}

	return false
}

///////////////////////////////////////////////////////////////////////////////

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

	// part 2
	safeCountDampened := 0
	for _, report := range reports {
		if isReportSafeDampened(report) {
			safeCountDampened++
		}
	}
	fmt.Println("2.1:", safeCountDampened)

}
