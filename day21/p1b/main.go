package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/graph"
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

var (
	numPadGraph = getNumPadGraph()
	dirPadGraph = getDirPadGraph()
)

// 177830 too high
// 183830 too high

// 182112 too high
// 182810 not right
func main() {

	lines, _ := file.GetLines("../test.txt")

	sum := 0

	//for _, line := range lines {
	//	seqs := getSequencesForSequence(numPadGraph, line)
	//	for _, seq := range seqs {
	//		seqs2 := getSequencesForSequence(dirPadGraph, seq)
	//		fmt.Println(seq, seqs2)
	//	}
	//	break
	//}

	for _, line := range lines {
		seq := "A" + line
		var sequences []string
		for i := 0; i < len(seq)-1; i++ {
			sequences = combineSequences(sequences, getPathSeq(seq[i:i+2], 2))
		}

		fmt.Println(sequences)
		fmt.Println(len(sequences), len(sequences[0]))
		break
	}

	fmt.Printf("\n\nAnswer: %d\n", sum)

}

type DepthSequence struct {
	s string
	d int
}

var (
	memDepthSequence = make(map[DepthSequence][]string)
)

func getPathSeq(sequence string, depth int) []string {
	if s, e := memDepthSequence[DepthSequence{sequence, depth}]; e {
		return s
	}
	if depth == 0 {
		return getPathsForSequence(numPadGraph, sequence)
	}
	prevSequences := getPathSeq(sequence, depth-1)
	var nextSequences []string
	for _, prevSeq := range prevSequences {
		var sequences []string
		seq := "A" + prevSeq
		for i := 0; i < len(seq)-1; i++ {
			sequences = combineSequences(sequences, getPathsForSequence(dirPadGraph, seq[i:i+2]))
		}
		if len(nextSequences) == 0 || len(sequences[0]) < len(nextSequences[0]) {
			nextSequences = sequences
		} else if len(sequences[0]) == len(nextSequences[0]) {
			nextSequences = append(nextSequences, sequences...)
		}
	}
	memDepthSequence[DepthSequence{sequence, depth}] = nextSequences
	return nextSequences
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

type SeqMem map[string][]string

func (sm SeqMem) Has(seq string) bool {
	if _, e := sm[seq]; e {
		return true
	}
	return false
}

type GraphSeqMem map[*graph.Graph]SeqMem

var (
	memSeq = make(GraphSeqMem)
)

func (gsm GraphSeqMem) Has(g *graph.Graph, seq string) bool {
	if _, e := gsm[g]; !e {
		gsm[g] = make(SeqMem)
		return false
	} else if gsm[g].Has(seq) {
		return true
	}
	return false
}

func getPathsForSequence(g *graph.Graph, seq string) []string {
	if memSeq.Has(g, seq) {
		return memSeq[g][seq]
	}

	var sequences []string

	type state struct {
		n   *graph.Node
		seq string
	}

	for i := 0; i < len(seq)-1; i++ {

		queue := make(common.Queue[state], 0, g.Len())

		start := g.GetNode(seq[i : i+1])
		goal := g.GetNode(seq[i+1 : i+2])

		queue.Enqueue(state{n: start, seq: ""})

		for !queue.Empty() {
			cur := *(queue.Dequeue())
			if cur.n == goal {
				seqWithA := cur.seq + "A"
				if len(sequences) == 0 || len(seqWithA) == len(sequences[0]) {
					sequences = append(sequences, seqWithA)
				} else if len(seqWithA) < len(sequences[0]) {
					sequences = []string{seqWithA}
				} else {
					break
				}
			} else {
				edges := cur.n.GetTraversableEdges()
				for _, edge := range edges {
					dest := edge.GetDestination()
					destSeq := cur.seq + edge.GetProperty("dir").(string)
					if len(sequences) == 0 || len(destSeq) <= len(sequences[0]) {
						queue.Enqueue(state{n: dest, seq: destSeq})
					}
				}
			}

		}
	}

	return sequences
}

func combineSequences(previousSequences, newSequences []string) []string {
	if len(previousSequences) == 0 {
		return newSequences
	}
	product := common.CartesianProduct([][]string{previousSequences, newSequences})
	sequences := make([]string, len(product))
	for i, p := range product {
		sequences[i] = strings.Join(p, "")
	}
	return sequences
}

func getSequencesForSequence(g *graph.Graph, seq string) []string {
	var sequences []string

	seqFromA := "A" + seq

	for i := 0; i < len(seqFromA)-1; i++ {
		paths := getPathsForSequence(g, seqFromA[i:i+2])
		if len(sequences) == 0 {
			sequences = paths
		} else {
			sequences = combineSequences(sequences, paths)
		}
	}

	return sequences
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
