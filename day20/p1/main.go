package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"github.com/mbordner/aoc2024/common/graph/djikstra"
	"sort"
)

// 5520 too high
func main() {
	s, e, g := getData("../data.txt")

	grid := getGrid(g)

	shortestPaths := djikstra.GenerateShortestPaths(g, s)
	path, c := shortestPaths.GetShortestPath(e)

	// does not contain start
	pathPosContainer := make(common.PosContainer)
	pathPositions := make(common.Positions, len(path))
	startPosition := s.GetID().(common.Pos)

	prev := startPosition
	for i, n := range path {
		p := n.GetID().(common.Pos)
		pathPositions[i] = p
		pathPosContainer[p] = true
		r := rune('*')
		if p.Y == prev.Y {
			if p.X > prev.X {
				r = rune('>')
			} else {
				r = rune('<')
			}
		} else if p.X == prev.X {
			if p.Y > prev.Y {
				r = rune('v')
			} else {
				r = rune('^')
			}
		}
		if n != e {
			grid[p.Y][p.X] = byte(r)
		}
		prev = p
	}

	grid.Print()

	cost := int(c)

	fmt.Println(s, e, g.Len())
	fmt.Println("shortest path cost:", cost)

	cheats := calculateCheats(startPosition, pathPositions, pathPosContainer)

	saves := make([]int, 0, len(cheats))
	for s := range cheats {
		saves = append(saves, s)
	}
	sort.Ints(saves)

	count := 0
	for _, s := range saves {
		if s <= 100 {
			count += len(cheats[s])
		}
		fmt.Printf("There are %d cheats that save %d picoseconds\n", len(cheats[s]), s)
	}

	fmt.Println("number of cheats that would save you at least 100 picoseconds:", count)

}

type Cheat struct {
	P1 common.Pos
	P2 common.Pos
}

type Cheats map[int][]Cheat

func (cs Cheats) Add(c Cheat, reduction int) {
	if _, e := cs[reduction]; e {
		cs[reduction] = append(cs[reduction], c)
	} else {
		cs[reduction] = []Cheat{c}
	}
}

func calculateCheats(sp common.Pos, path common.Positions, pc common.PosContainer) Cheats {
	pathWithStart := append(common.Positions{sp}, path...)

	cheats := make(Cheats)

	step := make(map[common.Pos]int)
	for i, p := range pathWithStart {
		step[p] = i
	}

	for _, p := range pathWithStart {
		possibleCheats := make([]Cheat, 0, 16)
		for y1 := p.Y - 1; y1 <= p.Y+1; y1++ {
			for x1 := p.X - 1; x1 <= p.X+1; x1++ {
				if (x1 == p.X && y1 != p.Y) || (y1 == p.Y && x1 != p.X) {
					for y2 := y1 - 1; y2 <= y1+1; y2++ {
						for x2 := x1 - 1; x2 <= x1+1; x2++ {
							if (y2 == y1 && x2 != x1) || (x2 == x1 && y2 != y1) {
								p1 := common.Pos{Y: y1, X: x1}
								p2 := common.Pos{Y: y2, X: x2}
								if !pc.Has(p1) && pc.Has(p2) {
									if step[p2] > step[p] {
										possibleCheats = append(possibleCheats, Cheat{P1: p1, P2: p2})
									}
								}
							}
						}
					}
				}
			}
		}

		for _, pc := range possibleCheats {
			reduction := step[pc.P2] - step[p] - 2
			if reduction > 0 {
				cheats.Add(pc, reduction)
			}
		}
	}

	return cheats
}

func getGrid(g *graph.Graph) common.Grid {
	chars := make(map[common.Pos]rune)

	ps := make(common.Positions, 0, g.Len())

	for _, n := range g.GetNodes() {
		p := n.GetID().(common.Pos)
		r := n.GetProperty("char").(rune)
		chars[p] = r
		ps = append(ps, p)
	}

	min, max := ps.Extents()

	grid := make(common.Grid, max.Y-min.Y+1)
	for y := range grid {
		grid[y] = make([]byte, max.X-min.X+1)
		for x := range grid[y] {
			p := common.Pos{Y: y, X: x}
			grid[y][x] = byte(chars[p])
		}
	}

	return grid
}

func getData(f string) (*graph.Node, *graph.Node, *graph.Graph) {
	lines, _ := file.GetLines(f)

	g := graph.NewGraph()

	var start, end *graph.Node

	for y, line := range lines {
		for x, r := range line {

			p := common.Pos{Y: y, X: x}
			n := g.CreateNode(p)

			n.AddProperty("char", rune(r))

			switch r {
			case '#':
				n.SetTraversable(false)
			case '.':
			case 'S':
				start = n
			case 'E':
				end = n
			}

		}
	}

	for y := 1; y < len(lines)-1; y++ {
		for x := 1; x < len(lines[y])-1; x++ {
			p := common.Pos{Y: y, X: x}
			n := g.GetNode(p)

			if y > 1 {
				o := g.GetNode(common.Pos{X: p.X, Y: p.Y - 1})
				n.AddEdge(o, 1)
				o.AddEdge(n, 1)
			}
			if x > 1 {
				o := g.GetNode(common.Pos{X: p.X - 1, Y: p.Y})
				n.AddEdge(o, 1)
				o.AddEdge(n, 1)
			}
			if x < len(lines[y])-1 {
				o := g.GetNode(common.Pos{X: p.X + 1, Y: p.Y})
				n.AddEdge(o, 1)
				o.AddEdge(n, 1)
			}
			if y < len(lines[y])-1 {
				o := g.GetNode(common.Pos{X: p.X, Y: p.Y + 1})
				n.AddEdge(o, 1)
				o.AddEdge(n, 1)
			}
		}
	}

	return start, end, g
}
