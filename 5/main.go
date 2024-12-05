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
	Before, After Page
}

type Page int

type Update []Page

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
			before, _ := strconv.Atoi(fields[0])
			after, _ := strconv.Atoi(fields[1])
			rules.Orderings = append(rules.Orderings, PageOrdering{Before: Page(before), After: Page(after)})
		}
	}

	for _, line := range strings.Split(rulesStr[separator+2:], "\n") {
		var update Update
		pages := strings.Split(line, ",")
		for _, page := range pages {
			pageNum, _ := strconv.Atoi(page)
			update = append(update, Page(pageNum))
		}
		rules.Updates = append(rules.Updates, update)
	}
	return rules
}

///////////////////////////////////////////////////////////////////////////////

func (r *Rules) isOrderValid(x, y Page) bool {
	for _, ordering := range r.Orderings {
		if ordering.Before == y && ordering.After == x {
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

func (r *Rules) repairUpdate(update Update) Update {
	// gotta keep re-sorting until done sorting
	// assumes no cycles in Orderings
	dirty := true
	for dirty {
		dirty = false
		for i := 0; i < len(update); i++ {
			pageI := update[i]
			for j := i + 1; j < len(update); j++ {
				// j is always after i
				pageJ := update[j]
				for _, ordering := range r.Orderings {
					if ordering.Before == pageJ && ordering.After == pageI {
						// move pageJ to where pageI was, shift rest over
						for k := j; k > i; k-- {
							update[k] = update[k-1]
						}
						update[i] = pageJ
						dirty = true // gotta keep re-sorting until done
					}
				}
			}
		}
	}
	return update
}

func (r *Rules) findAndRepairUpdates() []Update {
	var repaired []Update

	for _, update := range r.Updates {
		if !r.isUpdateCorrect(update) {
			fmt.Printf("%v\n", update)
			rp := r.repairUpdate(update)
			repaired = append(repaired, rp)
			fmt.Printf("%v\n", rp)
		}
	}
	return repaired
}

///////////////////////////////////////////////////////////////////////////////

func sumUpdateMiddlePages(updates []Update) int {
	sumMiddles := 0
	for _, update := range updates {
		// extract center page
		middlePage := update[len(update)/2]
		sumMiddles += int(middlePage)
	}
	return sumMiddles
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
	sumMiddles := sumUpdateMiddlePages(correctUpdates)
	fmt.Println("5.1:", sumMiddles)

	// part 2
	repairedUpdates := rules.findAndRepairUpdates()
	sumMiddles = sumUpdateMiddlePages(repairedUpdates)
	fmt.Println("5.2:", sumMiddles)
}
