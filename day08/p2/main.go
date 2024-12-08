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

func (p Pos) antinodes(grid Grid, o Pos) Positions {
	antis := Positions{p}
	var a, b = p, o
	for {
		c := Pos{X: a.X - b.X + a.X, Y: a.Y - b.Y + a.Y}
		if grid.containsPos(c) {
			antis = append(antis, c)
			a, b = c, a
		} else {
			break
		}
	}
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

func (antinodes Antinodes) add(ant byte, anti Pos) {
	if ants, e := antinodes[anti]; e {
		antinodes[anti] = append(ants, ant)
	} else {
		antinodes[anti] = []byte{ant}
	}
}

func (antennas Antennas) add(ant byte, p Pos) {
	if antenna, exists := antennas[ant]; exists {
		antennas[ant] = append(antenna, p)
	} else {
		antennas[ant] = Positions{p}
	}
}

func main() {
	debug := false

	grid := getData("../data.txt")
	antennas := getAntennas(grid)
	antinodes := make(Antinodes)

	for ant, locs := range antennas {
		locPairs := common.GetPairSets[Pos](locs)
		for _, pair := range locPairs {
			antis := pair[0].antinodes(grid, pair[1])
			antis = append(antis, pair[1].antinodes(grid, pair[0])...)
			for _, anti := range antis {
				antinodes.add(ant, anti)
			}
		}
	}

	if debug {
		for p := range antinodes {
			grid[p.Y][p.X] = '#'
		}

		for _, line := range grid {
			fmt.Println(string(line))
		}
	}

	fmt.Println(len(antinodes))
}

func getAntennas(grid Grid) Antennas {
	antennas := make(Antennas)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if reAntenna.MatchString(string(grid[y][x])) {
				antennas.add(grid[y][x], Pos{X: x, Y: y})
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
