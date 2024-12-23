package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"slices"
	"sort"
	"strings"
)

type ConnectedNodeList []*graph.Node

func (l *ConnectedNodeList) Add(n *graph.Node) {
	if !slices.Contains(*l, n) {
		*l = append(*l, n)
	}
}

func (l *ConnectedNodeList) String() string {
	ids := make([]string, len(*l))
	for i, n := range *l {
		ids[i] = n.GetID().(string)
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

func GetConnectedNodeListForNode(n *graph.Node) ConnectedNodeList {
	edges := n.GetTraversableEdges()
	l := make(ConnectedNodeList, 0, len(edges))
	for _, edge := range edges {
		l.Add(edge.GetDestination())
	}
	return l
}

// ce,ch,fc,he,hk,ji,jm,ns,qj,vv,yl,zy   <- not right
// do,dp,dq,es,ij,kh,lb,mj,ob,qw,sf,zi
func main() {
	g := getGraph("../data.txt")
	fmt.Println("node count", g.Len())
	maxCnl := ConnectedNodeList{}

	for _, n := range g.GetNodes() {
		cnl := GetConnectedNodeListForNode(n)
		fmt.Println(n.GetID().(string), len(cnl), cnl.String())

	nextList:
		for len(cnl) > 0 {
			for _, o := range cnl {
				ocnl := GetConnectedNodeListForNode(o)
				checkList := ConnectedNodeList(common.FilterArray(cnl, ConnectedNodeList{o}))
				for _, cn := range checkList {
					if !slices.Contains(ocnl, cn) {
						cnl = checkList
						continue nextList
					}
				}
			}
			break
		}

		if len(cnl) > 0 {
			cnl = append(ConnectedNodeList{n}, cnl...)
			if len(cnl) > len(maxCnl) {
				maxCnl = cnl
			}
		}
	}

	fmt.Println(maxCnl.String())
	fmt.Println(len(maxCnl))

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
