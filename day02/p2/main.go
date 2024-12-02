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

func checkDecreasing(l []int, errors int) bool {
	val := l[0]
	// check if decreasing
	for i := 1; i < len(l); i++ {
		if l[i] >= val {
			if errors == 0 {
				if i > 1 && l[i] < l[i-2] {
					val = l[i-2]
					i--
				}
				errors++
				continue
			}
			return false
		} else if val-l[i] > 3 {
			if errors == 0 {
				errors++
				continue
			}
			return false
		} else {
			val = l[i]
		}
	}
	return true
}

func checkIncreasing(l []int, errors int) bool {
	val := l[0]
	// check if increasing
	for i := 1; i < len(l); i++ {
		if l[i] <= val {
			if errors == 0 {
				if i > 1 && l[i] > l[i-2] {
					val = l[i-2]
					i--
				}
				errors++
				continue
			}
			return false
		} else if l[i]-val > 3 {
			if errors == 0 {
				errors++
				continue
			}
			return false
		} else {
			val = l[i]
		}
	}
	return true
}

func check(l []int) bool {
	return checkIncreasing(l, 0) || checkDecreasing(l, 0) ||
		checkIncreasing(l[1:], 1) || checkDecreasing(l[1:], 1)
}

// 382
// 397
func main() {
	lines, _ := file.GetLines("../data.txt")
	count := 0
	for _, line := range lines {
		if check(getLevels(line)) {
			fmt.Println("+", line)
			count++
		} else {
			fmt.Println("-", line)
		}
	}
	fmt.Println(count)
}
