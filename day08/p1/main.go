package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
)

var (
	reAntenna = regexp.MustCompile(`[0-9|A-Z|a-z]`)
)

type Pos struct {
	X int
	Y int
}

type Grid [][]byte
type Positions []Pos

type Antennas map[byte]Positions
type Antinodes map[Pos][]byte

func (p Pos) antinodes(o Pos) Positions {
	antis := make(Positions, 2)
	antis[0] = Pos{X: o.X - p.X + o.X, Y: o.Y - p.Y + o.Y}
	antis[1] = Pos{X: p.X - o.X + p.X, Y: p.Y - o.Y + p.Y}
	return antis
}

func (g Grid) inRange(x, y int) bool {
	if y < 0 || y >= len(g) || x < 0 || x >= len(g[y]) {
		return false
	}
	return true
}

func (g Grid) containsPos(p Pos) bool {
	return g.inRange(p.X, p.Y)
}

func main() {
	grid := getData("../data.txt")
	antennas := getAntennas(grid)
	antinodes := make(Antinodes)

	for ant, locs := range antennas {
		locPairs := common.GetPairSets[Pos](locs)
		for _, pair := range locPairs {
			antis := pair[0].antinodes(pair[1])
			for _, anti := range antis {
				if grid.containsPos(anti) {
					if ants, e := antinodes[anti]; e {
						antinodes[anti] = append(ants, ant)
					} else {
						antinodes[anti] = []byte{ant}
					}
				}
			}
		}
	}

	fmt.Println(len(antinodes))
}

func getAntennas(grid Grid) Antennas {
	antennas := make(Antennas)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if reAntenna.MatchString(string(grid[y][x])) {
				p := Pos{X: x, Y: y}
				if antenna, exists := antennas[grid[y][x]]; exists {
					antennas[grid[y][x]] = append(antenna, p)
				} else {
					antennas[grid[y][x]] = Positions{p}
				}
			}
		}
	}
	return antennas
}

func getData(f string) Grid {
	lines, _ := file.GetLines(f)
	grid := make(Grid, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return grid
}
