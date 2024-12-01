package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
)

var (
	reIds = regexp.MustCompile(`(\d+)\s+(\d+)`)
)

func parse(a, b string) (int64, int64) {
	x, _ := strconv.ParseInt(a, 10, 64)
	y, _ := strconv.ParseInt(b, 10, 64)
	return x, y
}

func getCount(a int64, counts map[int64]int) int {
	if c, e := counts[a]; e {
		return c
	}
	return 0
}

func main() {
	lines, _ := file.GetLines("../data.txt")

	lists := make([][]int64, 2, 2)
	lists[0] = make([]int64, len(lines), len(lines))
	lists[1] = make([]int64, len(lines), len(lines))

	counts := make(map[int64]int)

	for i, line := range lines {
		matches := reIds.FindAllStringSubmatch(line, -1)
		a, b := parse(matches[0][1], matches[0][2])
		lists[0][i], lists[1][i] = a, b
		if c, e := counts[b]; e {
			counts[b] = c + 1
		} else {
			counts[b] = 1
		}
	}

	var score int64
	for i := 0; i < len(lines); i++ {
		a := lists[0][i]
		score += a * int64(getCount(a, counts))
	}

	fmt.Println(score)
}
