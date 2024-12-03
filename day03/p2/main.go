package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
	"strings"
)

var (
	reMul = regexp.MustCompile(`mul\((\d+),(\d+)\)`)
)

func getVal(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func main() {

	var sum int64
	content, _ := file.GetContent("../data.txt")

	splits := strings.Split(string(content), `don't()`)
	updated := make([]byte, 0, len(content))

	updated = append(updated, splits[0]...)

	for i := 1; i < len(splits); i++ {
		doStart := strings.Index(splits[i], `do()`)
		if doStart >= 0 {
			updated = append(updated, splits[i][doStart:]...)
		}
	}

	matches := reMul.FindAllStringSubmatch(string(updated), -1)
	for _, match := range matches {
		sum += getVal(match[1]) * getVal(match[2])
	}

	fmt.Println(sum)
}
