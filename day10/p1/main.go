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
	for i := -1; i <= 1; i++ {
		if i != 0 {
			r := Pos{X: p.X + i, Y: p.Y}
			c := Pos{X: p.X, Y: p.Y + i}
			if grid.containsPos(r) {
				ps = append(ps, r)
			}
			if grid.containsPos(c) {
				ps = append(ps, c)
			}
		}
	}
	return ps
}

func trails(grid Grid, pos Pos) []Positions {
	if pos.val(grid) == 9 {
		return []Positions{{pos}}
	}
	ns := pos.getNeighborPositionsInGrid(grid)
	ps := make([]Positions, 0, len(ns))
	for _, n := range ns {
		nv := n.val(grid)
		pv := pos.val(grid)
		if nv == pv+1 {
			rss := trails(grid, n)
			for _, rs := range rss {
				ps = append(ps, append(Positions{pos}, rs...))
			}
		}
	}
	return ps
}

func main() {
	grid, posMap := getData("../data.txt")

	sum := 0
	for _, p := range posMap[0] {
		results := trails(grid, p)
		goalRouteCount := make(map[Pos]int)
		for _, ps := range results {
			goal := ps[len(ps)-1]
			if c, e := goalRouteCount[goal]; e {
				goalRouteCount[goal] = c + 1
			} else {
				goalRouteCount[goal] = 1
			}
		}
		sum += len(goalRouteCount)
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
