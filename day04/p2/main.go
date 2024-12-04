package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strings"
)

func check2(grid []string, x1, y1, x2, y2 int) bool {
	set := `MS`
	if strings.Contains(set, string(grid[y1][x1])) &&
		strings.Contains(set, string(grid[y2][x2])) &&
		grid[y1][x1] != grid[y2][x2] {
		return true
	}
	return false
}

func check(grid []string, y int, x int) bool {
	if x < 1 || x > len(grid[y])-2 || y < 1 || y > len(grid)-2 {
		return false
	}
	if check2(grid, x-1, y-1, x+1, y+1) && check2(grid, x+1, y-1, x-1, y+1) {
		return true
	}
	return false
}

func main() {

	count := 0
	grid, _ := file.GetLines("../data.txt")

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == 'A' {
				if check(grid, y, x) {
					count++
				}
			}
		}
	}

	fmt.Println(count)
}
