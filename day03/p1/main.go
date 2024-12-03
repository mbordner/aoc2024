package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
)

var (
	reMul = regexp.MustCompile(`mul\((\d+),(\d+)\)`)
)

func getVal(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func main() {

	lines, _ := file.GetLines("../data.txt")

	var sum int64

	for _, line := range lines {
		matches := reMul.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			sum += getVal(match[1]) * getVal(match[2])
		}
	}

	fmt.Println(sum)

}
