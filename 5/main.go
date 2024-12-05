// https://adventofcode.com/2024/day/5
// go run 5/main.go 5/5.txt

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PageOrdering struct {
	A, B int
}

type Update []int

type Rules struct {
	Orderings []PageOrdering
	Updates   []Update
}

///////////////////////////////////////////////////////////////////////////////

func NewRules(rulesStr string) *Rules {
	rules := &Rules{}

	separator := strings.Index(rulesStr, "\n\n")
	for _, line := range strings.Split(rulesStr[:separator], "\n") {
		fields := strings.Split(line, "|")
		if len(fields) == 2 {
			a, _ := strconv.Atoi(fields[0])
			b, _ := strconv.Atoi(fields[1])
			rules.Orderings = append(rules.Orderings, PageOrdering{A: a, B: b})
		}
	}

	for _, line := range strings.Split(rulesStr[separator+2:], "\n") {
		var update Update
		pages := strings.Split(line, ",")
		for _, page := range pages {
			pageNum, _ := strconv.Atoi(page)
			update = append(update, pageNum)
		}
		rules.Updates = append(rules.Updates, update)
	}
	return rules
}

///////////////////////////////////////////////////////////////////////////////

func (r *Rules) isOrderValid(a, b int) bool {
	for _, ordering := range r.Orderings {
		if ordering.A == b && ordering.B == a {
			return false
		}
	}
	return true
}

func (r *Rules) isUpdateCorrect(update Update) bool {
	for i, page := range update {
		// check validity of page
		for j := i + 1; j < len(update); j++ {
			if !r.isOrderValid(page, update[j]) {
				return false
			}
		}
	}
	return true
}

func (r *Rules) findCorrectUpdates() []Update {
	var correct []Update
	for _, update := range r.Updates {
		if r.isUpdateCorrect(update) {
			correct = append(correct, update)
		}
	}
	return correct
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	rulesData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	rules := NewRules(string(rulesData))
	if rules == nil {
		fmt.Fprintf(os.Stderr, "Error creating Rules\n")
		os.Exit(1)
	}

	// part 1
	correctUpdates := rules.findCorrectUpdates()
	sumMiddles := 0
	for _, update := range correctUpdates {
		// extract center page
		middlePage := update[len(update)/2]
		sumMiddles += middlePage
	}
	fmt.Println("3.2:", sumMiddles)
}
