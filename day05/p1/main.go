package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
	"strings"
)

var (
	reRule = regexp.MustCompile(`(\d+)\|(\d+)`)
)

type Rules map[int][]int

type Update struct {
	pageOrder []int
	pages     map[int]int
}

func NewUpdate(values []int) *Update {
	u := new(Update)
	u.pageOrder = values
	u.pages = make(map[int]int)
	for i, v := range values {
		u.pages[v] = i
	}
	return u
}

func validate(rules Rules, update *Update) bool {
	for _, p := range update.pageOrder {
		if after, e := rules[p]; e {
			// this p has to be before all rule,
			beforePageIndex := update.pages[p]
			for _, afterPage := range after {
				if afterPageIndex, e2 := update.pages[afterPage]; e2 {
					if beforePageIndex > afterPageIndex {
						return false
					}
				}
			}
		}
	}
	return true
}

func middle(update *Update) int {
	return update.pageOrder[len(update.pageOrder)/2]
}

func main() {
	rules, updates := getData("../data.txt")

	sum := 0

	for _, update := range updates {
		if validate(rules, update) {
			sum += middle(update)
		}
	}

	fmt.Println(sum)
}

func getVal(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}

func getVals(s string) []int {
	tokens := strings.Split(s, ",")
	values := make([]int, len(tokens), len(tokens))
	for i, v := range tokens {
		values[i] = getVal(v)
	}
	return values
}

func getData(f string) (Rules, []*Update) {
	lines, _ := file.GetLines(f)
	var rules = make(Rules)
	var updates = make([]*Update, 0, len(lines))

	parsingRules := true
	for _, line := range lines {
		if line == "" {
			parsingRules = false
			continue
		}
		if parsingRules {
			matches := reRule.FindStringSubmatch(line)
			v1 := getVal(matches[1])
			v2 := getVal(matches[2])
			if _, e := rules[v1]; e {
				rules[v1] = append(rules[v1], v2)
			} else {
				rules[v1] = []int{v2}
			}

		} else {
			updates = append(updates, NewUpdate(getVals(line)))
		}
	}

	return rules, updates
}
