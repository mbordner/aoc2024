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

func (n NodeId) toPos() common.Pos {
	return common.Pos{X: n.x, Y: n.y}
}

func main() {
	startId, goalId, g := getData("../data.txt")

	goalNode := g.GetNode(goalId)
	startShortestPaths := djikstra.GenerateShortestPaths(g, g.GetNode(startId))
	path, cost := startShortestPaths.GetShortestPath(goalNode)
	fmt.Println("shortest path length:", len(path), "cost:", int(cost))

	endShortestPaths := djikstra.GenerateShortestPaths(g, goalNode)

	optimalTiles := make(map[common.Pos]bool)
	optimalTiles[startId.toPos()] = true
	optimalTiles[goalId.toPos()] = true
	for _, n := range path {
		optimalTiles[n.GetID().(NodeId).toPos()] = true
	}

	traversableNodes := g.GetTraversableNodes()
	for _, n := range traversableNodes {
		nId := n.GetID().(NodeId)
		if _, e := optimalTiles[nId.toPos()]; !e {
			sToNpath, sToNcost := startShortestPaths.GetShortestPath(n)
			if len(sToNpath) > 0 && int(sToNcost) < int(cost) {
				nToEpath, nToEcost := endShortestPaths.GetShortestPath(getOpNode(g, nId))
				if len(nToEpath) > 0 {
					if int(sToNcost)+int(nToEcost) == int(cost) {
						optimalTiles[nId.toPos()] = true
					}
				}
			}

		}
	}

	fmt.Println("number of positions on all optimal paths:", len(optimalTiles))
}

func getOpNode(g *graph.Graph, id NodeId) *graph.Node {
	opDir := Any
	if id.d == N {
		opDir = S
	} else if id.d == S {
		opDir = N
	} else if id.d == W {
		opDir = E
	} else if id.d == E {
		opDir = W
	}
	return g.GetNode(NodeId{x: id.x, y: id.y, d: opDir})
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

			opN := getOpNode(g, id)
			goalNode := g.GetNode(goalId)
			goalNode.AddEdge(opN, 1)

		}
		n.AddEdge(g.GetNode(nextId), 1)
	}

	return startId, goalId, g
}
