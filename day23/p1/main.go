package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"sort"
	"strings"
)

type FNT struct {
	from *graph.Node
	node *graph.Node
	to   *graph.Node
}

func (fnt FNT) String() string {
	from := fnt.from.GetID().(string)
	node := fnt.node.GetID().(string)
	to := fnt.to.GetID().(string)

	//if strings.Compare(to, from) < 0 {
	//	from, to = to, from
	//}

	return fmt.Sprintf("%s,%s,%s", from, node, to)
}

func main() {
	g := getGraph("../data.txt")

	nodes := g.GetNodes()

	connections := make(map[FNT]bool)

	for _, fromNode := range nodes {
		for _, fromEdge := range fromNode.GetEdges() {
			node := fromEdge.GetDestination()
			for _, toEdge := range node.GetEdges() {
				toNode := toEdge.GetDestination()
				if toNode != fromNode && toNode != node {
					for _, checkEdge := range toNode.GetEdges() {
						if checkEdge.GetDestination() == fromNode {
							connections[FNT{fromNode, node, toNode}] = true
						}
					}
				}
			}
		}
	}

	connectionSets := make(map[string]bool)
	for c := range connections {
		cs := strings.Split(c.String(), ",")
		sort.Strings(cs)
		connectionSets[strings.Join(cs, ",")] = true
	}

	count := 0
nextSet:
	for cs := range connectionSets {
		tokens := strings.Split(cs, ",")
		for _, token := range tokens {
			if token[0] == 't' {
				count++
				continue nextSet
			}
		}
	}

	fmt.Println(count)

}

func getGraph(f string) *graph.Graph {
	g := graph.NewGraph()

	lines, _ := file.GetLines(f)
	for _, line := range lines {
		tokens := strings.Split(line, "-")

		var n, o *graph.Node
		if n = g.GetNode(tokens[0]); n == nil {
			n = g.CreateNode(tokens[0])
		}
		if o = g.GetNode(tokens[1]); o == nil {
			o = g.CreateNode(tokens[1])
		}

		n.AddEdge(o, 1)
		o.AddEdge(n, 1)
	}

	return g
}
