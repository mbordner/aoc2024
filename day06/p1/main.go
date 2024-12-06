package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/geom"
)

func main() {
	count := 0
	grid, x, y := getData("../data.txt")

	dir := geom.North
	for {

		nx := x
		ny := y

		switch dir {
		case geom.North:
			ny--
		case geom.East:
			nx++
		case geom.South:
			ny++
		case geom.West:
			nx--
		}
		if grid[y][x] != 'G' {
			grid[y][x] = 'G'
			count++
		}

		if ny < 0 || ny == len(grid) || nx < 0 || nx == len(grid[y]) {
			break
		} else if grid[ny][nx] == '#' {
			switch dir {
			case geom.North:
				dir = geom.East
			case geom.East:
				dir = geom.South
			case geom.South:
				dir = geom.West
			case geom.West:
				dir = geom.North
			}
		} else {
			y = ny
			x = nx
		}
	}

	fmt.Println(count)
}

func getData(f string) ([][]byte, int, int) {
	lines, _ := file.GetLines(f)
	grid := make([][]byte, len(lines), len(lines))
	var y, x int

	found := false
	for j := 0; j < len(lines); j++ {
		grid[j] = []byte(lines[j])
		if !found {
			for i := 0; i < len(grid[j]); i++ {
				if grid[j][i] == '^' {
					found = true
					y = j
					x = i
					break
				}
			}
		}
	}

	return grid, x, y
}
