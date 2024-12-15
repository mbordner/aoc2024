package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
)

type Vector struct {
	X int64
	Y int64
}

type Button struct {
	Cost int64
	V    Vector
}

// max int64: 18446744073709551615

type Machine struct {
	A     Button
	B     Button
	Prize Vector
}

var (
	reButtonLine = regexp.MustCompile(`^Button\s+([A|B]):\s+X([+|-])(\d+),\s+Y([+|-])(\d+)$`)
	rePrizeLine  = regexp.MustCompile(`^Prize:\s+X=(\d+), Y=(\d+)$`)
)

type State struct {
	a int64
	b int64
	m *Machine
}

func (s State) cost() int64 {
	return s.a*s.m.A.Cost + s.b*s.m.B.Cost
}

func solve(m *Machine) []State {
	solutions := make([]State, 0, 1)

	// Cramer's Rule to solve system of linear equations with 2 varialbes
	// given A Vector ( ax, ay ) and B Vector ( bx, by )  and Prize Vector (px, py )
	//    A * ax + B * bx = px
	//    A * ay + B * by = py,   where A and B are our variables, number of button presses

	// determinate of 2 x 2 matrix
	// [ a   b ]   ==  a*d - c*b
	// [ c   d ]

	// get determinate
	// det [ ax  bx ]
	//     [ ay  by ]   ax*by  - ay*bx
	// deta [ px   bx ]
	//      [ py   by ]    px*by - py*bx
	// detb [ ax   px ]
	//      [ ay   py ]   ax*py - ay*px
	// then A = deta / det
	//     B = detb / det
	det := m.A.V.X*m.B.V.Y - m.A.V.Y*m.B.V.X
	if det != 0 { // if it's 0, it can't solve for A & B because of divide by 0

		solution := State{m: m}

		deta := m.Prize.X*m.B.V.Y - m.Prize.Y*m.B.V.X
		detb := m.A.V.X*m.Prize.Y - m.A.V.Y*m.Prize.X

		solution.a = deta / det
		solution.b = detb / det

		if solution.a >= 0 && solution.b >= 0 {
			// check due to rounding
			if solution.a*m.A.V.X+solution.b*m.B.V.X == m.Prize.X &&
				solution.a*m.A.V.Y+solution.b*m.B.V.Y == m.Prize.Y {
				solutions = append(solutions, solution)
			}
		}

	}

	return solutions
}

func main() {
	machines := getData("../data.txt")

	tokens := int64(0)
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
			machine.Prize = Vector{X: x + int64(10000000000000), Y: y + int64(10000000000000)}
		}
	}

	if machine != nil {
		machines = append(machines, machine)
	}

	return machines
}

func getIntVal(sign, num string) int64 {
	if sign == "-" {
		num = sign + num
	}
	val, _ := strconv.ParseInt(num, 10, 64)
	return int64(val)
}
