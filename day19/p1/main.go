package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"sort"
	"strings"
)

type Towels map[string]bool

func (t Towels) Has(s string) bool {
	if _, e := t[s]; e {
		return true
	}
	return false
}

func (t Towels) String() string {
	ss := make([]string, 0, len(t))
	for k := range t {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return strings.Join(ss, ",")
}

type Stripes map[rune][]string

func (s Stripes) Has(r rune) bool {
	if _, e := s[r]; e {
		return true
	}
	return false
}

func candidates(towels Towels, stripes Stripes, pattern string) Towels {
	pTowels := make(Towels)
	for _, r := range pattern {
		if stripes.Has(r) {
			for _, s := range stripes[r] {
				pTowels[s] = true
			}
		} else {
			return Towels{}
		}
	}
	return pTowels
}

var (
	mem = make(map[string][]string)
)

func possible(towels Towels, stripes Stripes, pattern string) []string {
	if pattern == "" {
		return []string{""}
	}

	if ps, e := mem[pattern]; e {
		return ps
	}
	cs := candidates(towels, stripes, pattern)

	results := make(map[string]bool)

nextT:
	for towel := range cs {
		tokens := strings.Split(pattern, towel)
		if len(tokens) == 1 && tokens[0] == pattern {
			continue nextT
		} else {
			resultsMap := make(map[string]bool)
			left := possible(towels, stripes, tokens[0])
			if len(left) > 0 {
				for _, ls := range left {
					resultsMap[ls] = true
				}
				for i := 1; i < len(tokens); i++ {
					right := possible(towels, stripes, tokens[i])
					if len(right) == 0 {
						continue nextT
					} else {
						prevResultsMap := resultsMap
						resultsMap = make(map[string]bool)
						for ls := range prevResultsMap {
							for _, rs := range right {
								result := ls + "," + towel + "," + rs
								resultsMap[result] = true
							}
						}
					}
				}
			}
			for tr := range resultsMap {
				tr = strings.Join(common.Filter(strings.Split(tr, ","), ""), ",")
				if len(tr) > 0 {
					results[tr] = true
				}
			}
		}
	}

	combinations := make([]string, 0, len(results))
	for s := range results {
		combinations = append(combinations, s)
		break
	}

	mem[pattern] = combinations

	return combinations
}

func main() {
	towels, stripes, desiredPatterns := getData("../data.txt")

	count := 0
	for i, pattern := range desiredPatterns {
		fmt.Println("checking pattern (", i, ") ", pattern)
		results := possible(towels, stripes, pattern)
		if len(results) > 0 {
			count++
		}
	}

	fmt.Println(count, "designs possible")
}

func getData(f string) (Towels, Stripes, []string) {
	lines, _ := file.GetLines(f)
	towels := make(Towels)
	stripes := make(Stripes)

	tokens := strings.Split(lines[0], ", ")
	for _, token := range tokens {
		towels[token] = true
		for _, char := range token {
			if ss, e := stripes[char]; !e {
				stripes[char] = []string{token}
			} else {
				stripes[char] = append(ss, token)
			}
		}
	}

	desiredPatterns := make([]string, 0, len(lines)-2)
	for _, line := range lines[2:] {
		desiredPatterns = append(desiredPatterns, line)
	}

	return towels, stripes, desiredPatterns
}
