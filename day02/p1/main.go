package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
)

var (
	reDigits = regexp.MustCompile(`(\d+)`)
)

func getLevels(s string) []int {
	matches := reDigits.FindAllStringSubmatch(s, -1)
	levels := make([]int, len(matches), len(matches))
	for i := 0; i < len(matches); i++ {
		val, _ := strconv.ParseInt(matches[i][0], 10, 64)
		levels[i] = int(val)
	}
	return levels
}

func checkDecreasing(l []int) bool {
	val := l[0]
	// check if decreasing
	for i := 1; i < len(l); i++ {
		if l[i] >= val {
			return false
		} else if val-l[i] > 3 {
			return false
		} else {
			val = l[i]
		}
	}
	return true
}

func checkIncreasing(l []int) bool {
	val := l[0]
	// check if increasing
	for i := 1; i < len(l); i++ {
		if l[i] <= val {
			return false
		} else if l[i]-val > 3 {
			return false
		} else {
			val = l[i]
		}
	}
	return true
}

func check(l []int) bool {
	return checkIncreasing(l) || checkDecreasing(l)
}

func main() {
	lines, _ := file.GetLines("../data.txt")
	count := 0
	for _, line := range lines {
		if check(getLevels(line)) {
			//fmt.Println(line)
			count++
		}
	}
	fmt.Println(count)
}
