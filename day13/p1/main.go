package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"sort"
	"strconv"
)

type Vector struct {
	X int
	Y int
}

type Button struct {
	Cost int
	V    Vector
}

type Machine struct {
	A     Button
	B     Button
	Prize Vector
}

var (
	reButtonLine = regexp.MustCompile(`^Button\s+([A|B]):\s+X([+|-])(\d+),\s+Y([+|-])(\d+)$`)
	rePrizeLine  = regexp.MustCompile(`^Prize:\s+X=(\d+), Y=(\d+)$`)
)

const (
	maxPresses = 100
)

type State struct {
	a int
	b int
	m *Machine
}

func (s State) cost() int {
	return s.a*s.m.A.Cost + s.b*s.m.B.Cost
}

func (s State) checkGoal() int {
	if s.a > maxPresses || s.b > maxPresses {
		return 1
	}
	xa := s.a * s.m.A.V.X
	xb := s.b * s.m.B.V.X
	x := xa + xb
	if x > s.m.Prize.X {
		return 1
	}
	ya := s.a * s.m.A.V.Y
	yb := s.b * s.m.B.V.Y
	y := ya + yb
	if y > s.m.Prize.Y {
		return 1
	}
	if x == s.m.Prize.X && y == s.m.Prize.Y {
		return 0
	}
	return -1
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

func solve(m *Machine) []State {
	solutions := make([]State, 0, 100)

	queue := make(Queue, 0, 100)
	initial := State{a: 0, b: 0, m: m}
	queue.enqueue(initial)
	visited := make(Visited)
	visited[initial] = true

	for !queue.empty() {
		cur := queue.dequeue()
		g := cur.checkGoal()
		if g == 0 {
			solutions = append(solutions, *cur)
		} else if g < 0 {
			nextA := State{a: cur.a + 1, b: cur.b, m: m}
			if !visited.has(nextA) {
				visited[nextA] = true
				queue.enqueue(nextA)
			}
			nextB := State{a: cur.a, b: cur.b + 1, m: m}
			if !visited.has(nextB) {
				visited[nextB] = true
				queue.enqueue(nextB)
			}
		}
	}

	sort.Slice(solutions, func(i, j int) bool {
		if solutions[i].cost() < solutions[j].cost() {
			return true
		}
		return false
	})

	return solutions
}

func main() {
	machines := getData("../data.txt")

	tokens := 0
	for _, m := range machines {
		solutions := solve(m)
		if len(solutions) > 0 {
			tokens += solutions[0].cost()
		}
	}

	fmt.Println(tokens)
}

func getData(f string) []*Machine {
	lines, _ := file.GetLines(f)
	machines := make([]*Machine, 0, len(lines)/4+1)

	var machine *Machine
	for _, line := range lines {
		if line == "" {
			machines = append(machines, machine)
			machine = nil
			continue
		}
		if machine == nil {
			machine = &Machine{A: Button{Cost: 3}, B: Button{Cost: 1}}
		}
		if reButtonLine.MatchString(line) {
			matches := reButtonLine.FindStringSubmatch(line)
			button := matches[1]
			xOffset := getIntVal(matches[2], matches[3])
			yOffset := getIntVal(matches[4], matches[5])
			v := Vector{X: xOffset, Y: yOffset}
			switch button {
			case "A":
				machine.A.V = v
			case "B":
				machine.B.V = v
			}
		} else if rePrizeLine.MatchString(line) {
			matches := rePrizeLine.FindStringSubmatch(line)
			x := getIntVal("+", matches[1])
			y := getIntVal("+", matches[2])
			machine.Prize = Vector{X: x, Y: y}
		}
	}

	if machine != nil {
		machines = append(machines, machine)
	}

	return machines
}

func getIntVal(sign, num string) int {
	if sign == "-" {
		num = sign + num
	}
	val, _ := strconv.ParseInt(num, 10, 32)
	return int(val)
}
