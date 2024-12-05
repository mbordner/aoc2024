package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/array/ints"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	reRule = regexp.MustCompile(`(\d+)\|(\d+)`)
)

// Dependencies implements sort.Interface for sort.Sort(sort.Interface)
// first index is the page #, and all following values are page #s this page should be before
// it also represents the current page number list, as it is a list of page numbers with their dependent page numbers
type Dependencies [][]int

func (d Dependencies) Len() int {
	return len(d)
}

func (d Dependencies) Less(i, j int) bool {
	return ints.Contains(d[i][1:], d[j][0])
}

func (d Dependencies) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Rules is a simple type to collect all pages that should follow a given page
type Rules map[int][]int

// Update represents a page number list, can be valid, fixed so that it is valid, can get page list, and middle page
type Update struct {
	dependencies Dependencies
	pageIndexMap map[int]int
}

// NewUpdate builds an update with pages, and collects dependency information based on rules and other pages in its list
func NewUpdate(rules Rules, values []int) *Update {
	u := new(Update)
	u.dependencies = make(Dependencies, len(values), len(values))
	u.pageIndexMap = make(map[int]int)
	for i, v := range values {
		u.pageIndexMap[v] = i
	}
	for i := 0; i < len(values); i++ {
		dependencyValues := []int{values[i]}
		if after, e := rules[dependencyValues[0]]; e {
			for _, a := range after {
				if _, e2 := u.pageIndexMap[a]; e2 {
					dependencyValues = append(dependencyValues, a)
				}
			}
		}
		u.dependencies[i] = dependencyValues
	}
	return u
}

// pages returns page numbers in their current order for this update
func (u *Update) pages() []int {
	values := make([]int, len(u.dependencies), len(u.dependencies))
	for i, d := range u.dependencies {
		values[i] = d[0]
	}
	return values
}

// valid returns whether the order matches the rules
func (u *Update) valid() bool {
	for i := 0; i < len(u.dependencies)-2; i++ {
		if !u.dependencies.Less(i, i+1) {
			return false
		}
	}
	return true
}

// middle returns the middle page number
func (u *Update) middle() int {
	pages := u.pages()
	return pages[len(pages)/2]
}

// fixOrder corrects the page order based on the recorded dependencies
func (u *Update) fixOrder() {
	sort.Sort(u.dependencies)
}

func main() {
	updates := getData("../data.txt")

	sum := 0

	for _, update := range updates {
		if !update.valid() {
			update.fixOrder()
			sum += update.middle()
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

func getData(f string) []*Update {
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
			updates = append(updates, NewUpdate(rules, getVals(line)))
		}
	}

	return updates
}
