package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"github.com/mbordner/aoc2024/common/graph/djikstra"
)

type Dir int

const (
	Any Dir = iota
	N
	E
	S
	W
)

type NodeId struct {
	y int
	x int
	d Dir
}

func main() {
	startId, goalId, g := getData("../data.txt")

	shortestPaths := djikstra.GenerateShortestPaths(g, g.GetNode(startId))
	path, cost := shortestPaths.GetShortestPath(g.GetNode(goalId))
	fmt.Println(len(path), int(cost))
}

func getData(f string) (NodeId, NodeId, *graph.Graph) {

	lines, _ := file.GetLines(f)
	grid := common.ConvertGrid(lines)

	var startId NodeId
	var goalId NodeId

	g := graph.NewGraph()
	dirs := []Dir{N, E, S, W}

	for y, line := range grid {
		for x, c := range line {
			if c == 'S' || c == '.' || c == 'E' {
				for _, d := range dirs {
					nodeId := NodeId{x: x, y: y, d: d}
					g.CreateNode(nodeId)
				}
				if c == 'S' {
					startId = NodeId{x: x, y: y, d: E}
				} else if c == 'E' {
					goalId = NodeId{x: x, y: y, d: Any}
					g.CreateNode(goalId)
				}
			}
		}
	}

	for _, n := range g.GetNodes() {
		id := n.GetID().(NodeId)
		if id.x == goalId.x && id.y == goalId.y {
			continue
		}
		if id.d == N || id.d == S {
			n.AddEdge(g.GetNode(NodeId{x: id.x, y: id.y, d: E}), 1000)
			n.AddEdge(g.GetNode(NodeId{x: id.x, y: id.y, d: W}), 1000)
		} else if id.d == E || id.d == W {
			n.AddEdge(g.GetNode(NodeId{x: id.x, y: id.y, d: N}), 1000)
			n.AddEdge(g.GetNode(NodeId{x: id.x, y: id.y, d: S}), 1000)
		}
		nx, ny := id.x, id.y
		switch id.d {
		case Any:
		case N:
			ny--
		case E:
			nx++
		case S:
			ny++
		case W:
			nx--
		}
		nextId := NodeId{x: nx, y: ny, d: id.d}
		if nextId.x == goalId.x && nextId.y == goalId.y {
			nextId.d = Any
		}
		n.AddEdge(g.GetNode(nextId), 1)
	}

	return startId, goalId, g
}
