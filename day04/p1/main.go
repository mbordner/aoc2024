package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
)

// y, x == current x & y position in grid
// h == horizontal direction we're going in next position, negative for going left, positive for going right, 0 not changing columns
// y == vertical direction we're going in next position, negative for going up, positive for going down, 0 for not changing rows
// (we'll either be changing one or both of the row and column on our next search)
// p == current position in word that we matched
func next(grid []string, word []byte, y int, x int, h int, v int, p int) bool {
	np := p + 1
	if np == len(word) {
		return true
	}
	nx := x
	ny := y
	if h == -1 {
		if x == 0 {
			return false
		}
		nx--
	}
	if h == 1 {
		if x == len(grid[y])-1 {
			return false
		}
		nx++
	}
	if v == -1 {
		if y == 0 {
			return false
		}
		ny--
	}
	if v == 1 {
		if y == len(grid)-1 {
			return false
		}
		ny++
	}
	if grid[ny][nx] == word[np] {
		return next(grid, word, ny, nx, h, v, np)
	}
	return false
}

func main() {
	count := 0
	grid, _ := file.GetLines("../data.txt")

	word := []byte(`XMAS`)

	// left, top left, top, top right, right, bottom right, bottom, bottom left
	dirs := [][]int{{-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}}

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == word[0] {
				for _, d := range dirs {
					if next(grid, word, y, x, d[0], d[1], 0) {
						count++
					}
				}
			}
		}
	}

	fmt.Println(count)
}
