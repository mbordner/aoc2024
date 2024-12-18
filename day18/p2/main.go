package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
	"strings"
)

// cols and rows are the number of columns and rows in our space, start is 0,0 goal is cols,rows
// blockedPositions is an array of defined positions that would block a path, but when we check for shortest path, we aren't using all the blocked positions in the list
// we're just asking for the shortest path from start to goal if n number of blocked positions in the list existed
func getShortestPath(cols, rows int, blockedPositions common.Positions, n int) common.Positions {

	start := common.Pos{Y: 0, X: 0}
	goal := common.Pos{Y: rows, X: cols}

	queue := make(common.Queue[common.Pos], 0, cols*rows)
	visited := make(common.PosContainer)
	visited[start] = true

	blocked := make(common.PosContainer)
	for i := 0; i < n; i++ {
		blocked[blockedPositions[i]] = true
	}

	prev := make(common.PosLinker)
	queue.Enqueue(start)

	solution := common.Positions{}

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		if cur == goal {
			solution = common.Positions{goal}
			for p := prev[goal]; p != start; p = prev[p] {
				solution = append(common.Positions{p}, solution...)
			}
		} else {
			neighbors := make(common.Positions, 0, 4)
			for y := cur.Y - 1; y <= cur.Y+1; y++ {
				for x := cur.X - 1; x <= cur.X+1; x++ {
					// making sure this potential neighbor exists in bounds
					if y >= 0 && y <= rows && x >= 0 && x <= cols {
						// we're just constraining to current column and row, but skipping when matches current
						if (x != cur.X && y == cur.Y) || (y != cur.Y && x == cur.X) {
							p := common.Pos{Y: y, X: x}
							if !blocked.Has(p) {
								neighbors = append(neighbors, p)
							}
						}
					}
				}
			}
			for _, n := range neighbors {
				if !visited.Has(n) {
					visited[n] = true
					prev[n] = cur
					queue.Enqueue(n)
				}
			}
		}
	}

	return solution
}

type Params struct {
	cols      int
	rows      int
	numBlocks int
	file      string
}

func main() {

	params := []Params{{cols: 6, rows: 6, numBlocks: 12, file: "../test.txt"}, {cols: 70, rows: 70, numBlocks: 1024, file: "../data.txt"}}
	p := params[1]

	blockedPositions := getData(p.file)
	solution := getShortestPath(p.cols, p.rows, blockedPositions, p.numBlocks)
	fmt.Println(fmt.Sprintf("number of steps we could take on the shortest path if %d first blocked positions existed in the space: %d", p.numBlocks, len(solution)))

	i := p.numBlocks
	j := len(blockedPositions)

	// binary search through the blockedPositions array, starting at numBlocks since we already know we can find a path
	// if we used all the blocks we know we won't have a path, so we just have to search back through this ordered
	// list finding the first that will block our path
	for {
		m := (j-i)/2 + i
		//fmt.Println("checking ", m+1)
		s := getShortestPath(p.cols, p.rows, blockedPositions, m+1)
		if len(s) == 0 {
			if m-i == 0 {
				fmt.Println("position in the list of blocked positions where we would first have no shortest path solution:", m)
				fmt.Println(fmt.Sprintf("first %d blocked positions could find a path, but %d blocks would start to block a solution", m, m+1))
				fmt.Println(fmt.Sprintf("%d,%d", blockedPositions[m].X, blockedPositions[m].Y))
				break
			} else {
				j = m
			}
		} else {
			i = m + 1
		}
	}

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
