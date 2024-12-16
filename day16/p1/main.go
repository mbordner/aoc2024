package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"sort"
)

type Dir int

const (
	N Dir = iota
	E
	S
	W
)

type Reindeer struct {
	p common.Pos
	d Dir
}

type State struct {
	y int
	x int
	d Dir
}

type Queue []State

func (q *Queue) enqueue(s State) {
	*q = append(*q, s)
}

func (q *Queue) empty() bool {
	return len(*q) == 0
}

func (q *Queue) dequeue() *State {
	if !q.empty() {
		s := (*q)[0]
		*q = (*q)[1:]
		return &s
	}
	return nil
}

type Visited map[State]bool

func (v Visited) has(s State) bool {
	if b, e := v[s]; e {
		return b
	}
	return false
}

func getTurnOptions(d Dir) []Dir {
	if d == N || d == S {
		return []Dir{W, E}
	} else {
		return []Dir{N, S}
	}
}

func canTravelForward(grid common.Grid, dir Dir, x, y int) (int, int, bool) {
	nx, ny := x, y
	switch dir {
	case N:
		ny--
	case E:
		nx++
	case S:
		ny++
	case W:
		nx--
	}
	if grid[y][x] != '#' {
		return nx, ny, true
	}
	return nx, ny, false
}

type States []State

func (ss States) cost() int {
	cost := 0
	for i := 1; i < len(ss); i++ {
		if ss[i].d != ss[i-1].d {
			cost += 1000
		} else {
			cost += 1
		}
	}
	return cost
}

// 148508 too high

func main() {

	r, grid := getData("../data.txt")

	queue := make(Queue, 0, 100)
	initial := State{x: r.p.X, y: r.p.Y, d: r.d}
	queue.enqueue(initial)
	visited := make(Visited)
	visited[initial] = true
	prev := make(map[State]State)

	var solutions []States

	for !queue.empty() {
		cur := queue.dequeue()

		if grid[cur.y][cur.x] == 'E' {
			var solution States
			solution = States{*cur}
			for solution[0] != initial {
				solution = append(States{prev[solution[0]]}, solution...)
			}
			solutions = append(solutions, solution)
		} else {
			ns := make(States, 0, 3)
			for _, d := range getTurnOptions(cur.d) {
				ns = append(ns, State{x: cur.x, y: cur.y, d: d})
			}
			if nx, ny, can := canTravelForward(grid, cur.d, cur.x, cur.y); can {
				ns = append(ns, State{x: nx, y: ny, d: cur.d})
			}
			for _, s := range ns {
				if !visited.has(s) {
					visited[s] = true
					prev[s] = *cur
					queue.enqueue(s)
				}
			}
		}

	}

	sort.Slice(solutions, func(x, y int) bool {
		if solutions[x].cost() > solutions[y].cost() {
			return false
		}
		return true
	})

	ss := solutions[0]
	for i := 2; i < len(ss); i++ {
		if ss[i].x == ss[i-1].x && ss[i].x == ss[i-2].x &&
			ss[i].y == ss[i-1].y && ss[i].y == ss[i-2].y {
			fmt.Println("turn much?")
		}
	}

	for _, s := range solutions[0][1 : len(solutions[0])-1] {
		char := '^'
		if s.d == E {
			char = '>'
		} else if s.d == S {
			char = 'v'
		} else if s.d == W {
			char = '<'
		}
		grid[s.y][s.x] = byte(char)
	}
	for _, line := range grid {
		fmt.Println(string(line))
	}

	fmt.Println(solutions[0].cost())

}

func getData(f string) (Reindeer, common.Grid) {

	lines, _ := file.GetLines(f)
	grid := common.ConvertGrid(lines)

	r := Reindeer{d: E}

	for y, line := range grid {
		for x, c := range line {
			if c == 'S' {
				r.p.X, r.p.Y = x, y
			}
		}
	}

	return r, grid
}
