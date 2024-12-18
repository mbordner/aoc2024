package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
	"strings"
)

type Queue common.Positions

func (q *Queue) enqueue(s common.Pos) {
	*q = append(*q, s)
}

func (q *Queue) empty() bool {
	return len(*q) == 0
}

func (q *Queue) dequeue() *common.Pos {
	if !q.empty() {
		s := (*q)[0]
		*q = (*q)[1:]
		return &s
	}
	return nil
}

type PosContainer map[common.Pos]bool
type Previous map[common.Pos]common.Pos

func (v PosContainer) has(p common.Pos) bool {
	if b, e := v[p]; e {
		return b
	}
	return false
}

func getShortestPath(cols, rows int, walls common.Positions, n int) common.Positions {

	start := common.Pos{Y: 0, X: 0}
	goal := common.Pos{Y: rows, X: cols}

	queue := make(Queue, 0, cols*rows)
	visited := make(PosContainer)
	visited[start] = true

	blocked := make(PosContainer)
	for i := 0; i < n; i++ {
		blocked[walls[i]] = true
	}

	prev := make(Previous)
	queue.enqueue(start)

	solution := common.Positions{}

	for !queue.empty() {
		cur := *(queue.dequeue())
		if cur == goal {
			solution = common.Positions{goal}
			for p := prev[goal]; p != start; p = prev[p] {
				solution = append(common.Positions{p}, solution...)
			}
		} else {
			deltas := make(common.Positions, 0, 4)
			for y := cur.Y - 1; y <= cur.Y+1; y++ {
				for x := cur.X - 1; x <= cur.X+1; x++ {
					if y >= 0 && y <= rows && x >= 0 && x <= cols {
						if (x != cur.X && y == cur.Y) || (y != cur.Y && x == cur.X) {
							p := common.Pos{Y: y, X: x}
							if !blocked.has(p) {
								deltas = append(deltas, p)
							}
						}
					}
				}
			}
			for _, d := range deltas {
				if !visited.has(d) {
					visited[d] = true
					prev[d] = cur
					queue.enqueue(d)
				}
			}
		}
	}

	return solution
}

func main() {
	//walls := getData("../test.txt")
	//solution := getShortestPath(6, 6, walls, 12)

	walls := getData("../data.txt")
	solution := getShortestPath(70, 70, walls, 1024)

	fmt.Println(len(solution))

}

func getData(f string) common.Positions {
	lines, _ := file.GetLines(f)
	ps := make(common.Positions, 0, len(lines))
	for _, line := range lines {
		tokens := strings.Split(line, ",")
		p := common.Pos{X: getIntVal(tokens[0]), Y: getIntVal(tokens[1])}
		ps = append(ps, p)
	}
	return ps
}

func getIntVal(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
