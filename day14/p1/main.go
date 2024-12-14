package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
)

type Vector struct {
	X int
	Y int
}

type Robot struct {
	v Vector
	p Vector
}

type Robots []Robot

var (
	reRobot = regexp.MustCompile(`^p=(\d+),(\d+)\s+v=(-?\d+),(-?\d+)$`)
)

func (rs Robots) print(cols, rows int) {
	counts := make(map[Vector]int)
	for _, r := range rs {
		if c, e := counts[r.p]; e {
			counts[r.p] = c + 1
		} else {
			counts[r.p] = 1
		}
	}
	grid := make([][]byte, rows)
	for y := range grid {
		grid[y] = make([]byte, cols)
		for x := 0; x < len(grid[y]); x++ {
			grid[y][x] = '.'
		}
	}
	for p, c := range counts {
		grid[p.Y][p.X] = fmt.Sprintf("%d", c)[0]
	}
	for _, line := range grid {
		fmt.Println(string(line))
	}
}

func (rs Robots) move(cols, rows int) {

	for i, r := range rs {
		nx := r.p.X + r.v.X
		ny := r.p.Y + r.v.Y

		for nx < 0 {
			nx += cols
		}
		for ny < 0 {
			ny += rows
		}
		for nx >= cols {
			nx -= cols
		}
		for ny >= rows {
			ny -= rows
		}

		rs[i].p.X, rs[i].p.Y = nx, ny
	}

}

// 214701760 too low

// 223525120 high
func main() {
	cols, rows := 101, 103 // 11, 7   101, 103
	//cols, rows = 11, 7
	robots := getData("../data.txt")

	for i := 0; i < 100; i++ {

		robots.move(cols, rows)
		fmt.Println("=", i, "===================================================")
		robots.print(cols, rows)

	}

	xa := rows / 2
	ya := cols / 2

	quads := make([]int, 4)

	for _, r := range robots {
		if r.p.Y < xa && r.p.X > ya {
			quads[0]++
		} else if r.p.Y < xa && r.p.X < ya {
			quads[1]++
		} else if r.p.Y > xa && r.p.X < ya {
			quads[2]++
		} else if r.p.Y > xa && r.p.X > ya {
			quads[3]++
		}
	}

	val := quads[0]
	for _, q := range quads[1:] {
		val *= q
	}

	fmt.Println(val)

}

func getData(f string) Robots {
	lines, _ := file.GetLines(f)
	robots := make(Robots, len(lines))
	for i, line := range lines {
		matches := reRobot.FindStringSubmatch(line)
		px := getIntVal(matches[1])
		py := getIntVal(matches[2])
		vx := getIntVal(matches[3])
		vy := getIntVal(matches[4])
		robots[i] = Robot{p: Vector{X: px, Y: py}, v: Vector{X: vx, Y: vy}}
	}
	return robots
}

func getIntVal(num string) int {
	val, _ := strconv.ParseInt(num, 10, 32)
	return int(val)
}
