package main

import (
	"fmt"
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
	mem = make(map[string]int)
)

func possible(towels Towels, stripes Stripes, pattern string) int {
	if c, e := mem[pattern]; e {
		return c
	}
	result := 0
	cs := candidates(towels, stripes, pattern)
	for towel := range cs {
		if towel == pattern {
			result += 1
		} else if strings.HasPrefix(pattern, towel) {
			result += possible(towels, stripes, pattern[len(towel):])
		}
	}
	mem[pattern] = result
	return result
}

func main() {
	towels, stripes, desiredPatterns := getData("../data.txt")

	count := 0
	for i, pattern := range desiredPatterns {
		fmt.Println("checking pattern (", i, ") ", pattern)
		result := possible(towels, stripes, pattern)
		count += result
		fmt.Println("> designs: ", result)
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
