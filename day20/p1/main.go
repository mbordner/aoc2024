package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"github.com/mbordner/aoc2024/common/graph/djikstra"
)

type Cheat struct {
	p1        common.Pos
	p2        common.Pos
	reduction int
}

func main() {
	s, e, g := getData("../test.txt")

	grid := getGrid(g)

	shortestPaths := djikstra.GenerateShortestPaths(g, s)
	path, c := shortestPaths.GetShortestPath(e)

	// does not contain start
	pathPosContainer := make(common.PosContainer)
	pathPositions := make(common.Positions, len(path))
	startPosition := s.GetID().(common.Pos)

	prev := startPosition
	for i, n := range path[0 : len(path)-1] {
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
		grid[p.Y][p.X] = byte(r)
		prev = p
	}

	grid.Print()

	cost := int(c)

	fmt.Println(s, e, g.Len())
	fmt.Println("shortest path cost:", cost)

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
