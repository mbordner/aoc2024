package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"slices"
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

func sort(a, b int64) (int64, int64) {
	if a < b {
		return a, b
	}
	return b, a
}

func main() {
	lines, _ := file.GetLines("../data.txt")

	lists := make([][]int64, 2, 2)
	lists[0] = make([]int64, len(lines), len(lines))
	lists[1] = make([]int64, len(lines), len(lines))

	for i, line := range lines {
		matches := reIds.FindAllStringSubmatch(line, -1)
		lists[0][i], lists[1][i] = parse(matches[0][1], matches[0][2])
	}

	slices.Sort(lists[0])
	slices.Sort(lists[1])

	var diff int64
	for i := 0; i < len(lines); i++ {
		a, b := sort(lists[0][i], lists[1][i])
		diff += b - a
	}

	fmt.Println(diff)
}
