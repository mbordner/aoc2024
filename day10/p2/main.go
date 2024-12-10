package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
)

type Pos struct {
	X int
	Y int
}

type Positions []Pos
type Grid [][]int
type PosMap map[int]Positions

func (g Grid) inGrid(x, y int) bool {
	return y >= 0 && y < len(g) && x >= 0 && x < len(g[y])
}

func (g Grid) containsPos(p Pos) bool {
	return g.inGrid(p.X, p.Y)
}

func (p Pos) val(grid Grid) int {
	return grid[p.Y][p.X]
}

func (p Pos) getNeighborPositionsInGrid(grid Grid) Positions {
	ps := make(Positions, 0, 8)
	for y := p.Y - 1; y < p.Y+2; y++ {
		for x := p.X - 1; x < p.X+2; x++ {
			if (p.X != x && p.Y == y) || (p.X == x && p.Y != y) { // up, down, left, right
				o := Pos{X: x, Y: y}
				if grid.containsPos(o) {
					ps = append(ps, o)
				}
			}
		}
	}
	return ps
}

func trails(grid Grid, pos Pos) []Positions {
	if pos.val(grid) == 9 {
		return []Positions{{pos}}
	}
	neighbors := pos.getNeighborPositionsInGrid(grid)
	goalRoutes := make([]Positions, 0, len(neighbors))
	for _, n := range neighbors {
		if n.val(grid) == pos.val(grid)+1 {
			results := trails(grid, n)
			for _, result := range results {
				goalRoutes = append(goalRoutes, append(Positions{pos}, result...))
			}
		}
	}
	return goalRoutes
}

func main() {
	grid, posMap := getData("../data.txt")

	sum := 0
	for _, p := range posMap[0] {
		routesToGoal := trails(grid, p)
		goalRouteCounts := make(map[Pos]int)
		for _, ps := range routesToGoal {
			goal := ps[len(ps)-1]
			if c, e := goalRouteCounts[goal]; e {
				goalRouteCounts[goal] = c + 1
			} else {
				goalRouteCounts[goal] = 1
			}
		}
		sum += len(routesToGoal)
	}

	fmt.Println(sum)
}

func getData(f string) (Grid, PosMap) {
	lines, _ := file.GetLines(f)
	grid := make(Grid, len(lines))
	posMap := make(PosMap)
	for y, line := range lines {
		grid[y] = make([]int, len(line))
		for x := 0; x < len(line); x++ {
			pos := Pos{X: x, Y: y}
			val := getIntVal(string(line[x]))
			if ps, e := posMap[val]; e {
				posMap[val] = append(ps, pos)
			} else {
				posMap[val] = Positions{pos}
			}
			grid[y][x] = val
		}
	}
	return grid, posMap
}

func getIntVal(c string) int {
	val, err := strconv.ParseInt(c, 10, 32)
	if err != nil {
		return -1
	}
	return int(val)
}
