package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
	"github.com/mbordner/aoc2024/common/graph/djikstra"
	"regexp"
	"strconv"
	"strings"
)

var (
	reDigits = regexp.MustCompile(`^0*(\d+)`)
)

const (
	D = "v"
	L = "<"
	R = ">"
	U = "^"
)

type ShortestPathCache map[string]djikstra.ShortestPaths

func (spc ShortestPathCache) GetSP(g *graph.Graph, nodeId string) djikstra.ShortestPaths {
	if sp, e := spc[nodeId]; e {
		return sp
	} else {
		spc[nodeId] = djikstra.GenerateShortestPaths(g, g.GetNode(nodeId))
	}
	return spc[nodeId]
}

// 177830 too high
// 183830 too high

// 182112 too high
// 182810 not right
func main() {
	numPadGraph := getNumPadGraph()
	numPadSPC := make(ShortestPathCache)
	dirPadGraph := getDirPadGraph()
	dirPadSPC := make(ShortestPathCache)

	lines, _ := file.GetLines("../data.txt")

	sum := 0

	for _, line := range lines {
		val := getNumVal(line)

		fmt.Printf("\n>>>> processing line %s with value %d\n", line, val)
		seqNumPad := getSequenceForSequence(numPadGraph, numPadSPC, line)
		fmt.Printf("num pad seq: %s\n", pretty(seqNumPad))

		dirPadSeq := seqNumPad
		for i := 0; i < 2; i++ {
			dirPadSeq = getSequenceForSequence(dirPadGraph, dirPadSPC, dirPadSeq)
			fmt.Printf("dir pad seq %d: %s\n", i+1, pretty(dirPadSeq))
		}
		calculation := len(dirPadSeq) * val
		fmt.Printf("final length: %d, and calculation: %d\n", len(dirPadSeq), calculation)

		fmt.Println("npc len", len(numPadSPC), "dpc len", len(dirPadSPC))

		sum += calculation
	}

	fmt.Printf("\n\nAnswer: %d\n", sum)

}

func pretty(s string) string {
	return strings.Join(strings.Split(s, "A"), "A ")
}

func getNumVal(s string) int {
	if reDigits.MatchString(s) {
		match := reDigits.FindStringSubmatch(s)
		val, _ := strconv.Atoi(match[1])
		return val
	}
	return 0
}

func getSequenceForSequence(g *graph.Graph, spCache ShortestPathCache, seq string) string {
	sb := strings.Builder{}

	curId := "A"

	for i := 0; i < len(seq); i++ {

		nextId := string(seq[i])

		toSP := spCache.GetSP(g, curId)
		backSp := spCache.GetSP(g, nextId)

		toSB := strings.Builder{}
		backSB := strings.Builder{}

		_, toPathEdges, toCost := toSP.GetShortestPathWithEdges(g.GetNode(nextId))
		_, backPathEdges, backCost := backSp.GetShortestPathWithEdges(g.GetNode(curId))

		for _, e := range toPathEdges {
			toSB.WriteString(e.GetProperty("dir").(string))
		}

		for j := len(backPathEdges) - 1; j >= 0; j-- {
			e := backPathEdges[j]
			switch e.GetProperty("dir").(string) {
			case U:
				backSB.WriteString(D)
			case D:
				backSB.WriteString(U)
			case R:
				backSB.WriteString(L)
			case L:
				backSB.WriteString(R)
			}
		}

		if backCost < toCost && backSB.String() != toSB.String() {
			sb.WriteString(backSB.String())
		} else {
			sb.WriteString(toSB.String())
		}

		sb.WriteString("A")

		curId = nextId
	}

	return sb.String()
}

func getNumPadGraph() *graph.Graph {
	g := graph.NewGraph()

	bA := g.CreateNode("A")
	b0 := g.CreateNode("0")
	b1 := g.CreateNode("1")
	b2 := g.CreateNode("2")
	b3 := g.CreateNode("3")
	b4 := g.CreateNode("4")
	b5 := g.CreateNode("5")
	b6 := g.CreateNode("6")
	b7 := g.CreateNode("7")
	b8 := g.CreateNode("8")
	b9 := g.CreateNode("9")

	bA.AddEdge(b3, 1).AddProperty("dir", U)
	bA.AddEdge(b0, 1).AddProperty("dir", L)

	b0.AddEdge(bA, 1).AddProperty("dir", R)
	b0.AddEdge(b2, 1).AddProperty("dir", U)

	b1.AddEdge(b2, 1).AddProperty("dir", R)
	b1.AddEdge(b4, 1).AddProperty("dir", U)

	b2.AddEdge(b3, 1).AddProperty("dir", R)
	b2.AddEdge(b0, 1).AddProperty("dir", D)
	b2.AddEdge(b5, 1).AddProperty("dir", U)
	b2.AddEdge(b1, 1).AddProperty("dir", L)

	b3.AddEdge(bA, 1).AddProperty("dir", D)
	b3.AddEdge(b6, 1).AddProperty("dir", U)
	b3.AddEdge(b2, 1).AddProperty("dir", L)

	b4.AddEdge(b5, 1).AddProperty("dir", R)
	b4.AddEdge(b1, 1).AddProperty("dir", D)
	b4.AddEdge(b7, 1).AddProperty("dir", U)

	b5.AddEdge(b6, 1).AddProperty("dir", R)
	b5.AddEdge(b2, 1).AddProperty("dir", D)
	b5.AddEdge(b8, 1).AddProperty("dir", U)
	b5.AddEdge(b4, 1).AddProperty("dir", L)

	b6.AddEdge(b9, 1).AddProperty("dir", U)
	b6.AddEdge(b5, 1).AddProperty("dir", L)
	b6.AddEdge(b3, 1).AddProperty("dir", D)

	b7.AddEdge(b8, 1).AddProperty("dir", R)
	b7.AddEdge(b4, 1).AddProperty("dir", D)

	b8.AddEdge(b9, 1).AddProperty("dir", R)
	b8.AddEdge(b5, 1).AddProperty("dir", D)
	b8.AddEdge(b7, 1).AddProperty("dir", L)

	b9.AddEdge(b6, 1).AddProperty("dir", D)
	b9.AddEdge(b8, 1).AddProperty("dir", L)

	return g
}

func getDirPadGraph() *graph.Graph {
	g := graph.NewGraph()

	bA := g.CreateNode("A")
	bR := g.CreateNode(R)
	bU := g.CreateNode(U)
	bD := g.CreateNode(D)
	bL := g.CreateNode(L)

	bA.AddEdge(bR, 1).AddProperty("dir", D)
	bA.AddEdge(bU, 1).AddProperty("dir", L)

	bR.AddEdge(bA, 1).AddProperty("dir", U)
	bR.AddEdge(bD, 1).AddProperty("dir", L)

	bU.AddEdge(bA, 1).AddProperty("dir", R)
	bU.AddEdge(bD, 1).AddProperty("dir", D)

	bD.AddEdge(bU, 1).AddProperty("dir", U)
	bD.AddEdge(bR, 1).AddProperty("dir", R)
	bD.AddEdge(bL, 1).AddProperty("dir", L)

	bL.AddEdge(bD, 1).AddProperty("dir", R)

	return g
}
