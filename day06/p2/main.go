package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
)

type Dir int
type Grid [][]byte

const (
	NORTH Dir = iota
	EAST
	SOUTH
	WEST
)

type State struct {
	x int
	y int
	d Dir
}

type States []State

func next(grid Grid, cur State) *State {
	nx := cur.x
	ny := cur.y
	nd := cur.d

	switch nd {
	case NORTH:
		ny--
	case EAST:
		nx++
	case SOUTH:
		ny++
	case WEST:
		nx--
	}

	ns := State{x: nx, y: ny, d: nd}

	if ny < 0 || ny == len(grid) || nx < 0 || nx == len(grid[ny]) {
		return nil
	} else if grid[ny][nx] == '#' {
		nx, ny = cur.x, cur.y
		switch nd {
		case NORTH:
			nd = EAST
			nx++
		case EAST:
			nd = SOUTH
			ny++
		case SOUTH:
			nd = WEST
			nx--
		case WEST:
			nd = NORTH
			ny--
		}
		ns = State{x: cur.x, y: cur.y, d: nd}
	}

	return &ns
}

func main() {
	filename := "../data.txt"
	debug := false

	grid, start := getData(filename)

	path := States{start}

	for {
		cur := path[len(path)-1]

		if debug {
			if grid[cur.y][cur.x] != '^' && grid[cur.y][cur.x] != '+' {
				if cur.d == NORTH || cur.d == SOUTH {
					if grid[cur.y][cur.x] == '-' {
						grid[cur.y][cur.x] = '+'
					} else {
						grid[cur.y][cur.x] = '|'
					}
				} else {
					if grid[cur.y][cur.x] == '|' {
						grid[cur.y][cur.x] = '+'
					} else {
						grid[cur.y][cur.x] = '-'
					}
				}
			}
		}

		ns := next(grid, cur)
		if ns == nil {
			break
		}

		if debug {
			if cur.y == ns.y && cur.x == ns.x {
				grid[cur.y][cur.x] = '+'
			}
		}

		path = append(path, *ns)
	}

	visited := make(map[State]int)
	for _, p := range path {
		s := State{x: p.x, y: p.y}
		if c, e := visited[s]; e {
			visited[s] = c + 1
		} else {
			visited[s] = 1
		}
	}

	obstacles := make(map[State]int)

	for i := 1; i < len(path); i++ {
		obs := path[i]

		o := State{x: obs.x, y: obs.y}
		if _, e := obstacles[o]; e || (obs.y == start.y && obs.x == start.x) { // if we already know, or this is starting position skip
			continue
		}

		last := grid[obs.y][obs.x]
		grid[obs.y][obs.x] = '#'

		obsPath := States{path[0]}
		obsVisited := make(map[State]bool)
		obsVisited[path[0]] = true

		looping := false

		for {
			cur := obsPath[len(obsPath)-1]

			ns := next(grid, cur)
			if ns == nil {
				break
			}

			if _, e := obsVisited[*ns]; e {
				looping = true

				if debug {
					ogrid, _ := getData(filename)
					for _, s := range obsPath {
						char := '^'
						switch s.d {
						case NORTH:
							char = '^'
						case EAST:
							char = '>'
						case SOUTH:
							char = 'v'
						case WEST:
							char = '<'
						}
						ogrid[s.y][s.x] = byte(char)
					}
					ogrid[obs.y][obs.x] = 'O'
					ogrid[obsPath[0].y][obsPath[0].x] = 'S'
					ogrid[ns.y][ns.x] = '*'
					fmt.Println("=====================")
					for _, line := range ogrid {
						fmt.Println(string(line))
					}
					fmt.Println("=====================")
				}

				break
			} else {
				obsVisited[*ns] = true
			}
			obsPath = append(obsPath, *ns)

		}

		grid[obs.y][obs.x] = last

		if looping {
			o = State{x: obs.x, y: obs.y}
			if c, e := obstacles[o]; e {
				obstacles[o] = c + 1
			} else {
				obstacles[o] = 1
			}
		}
	}

	if debug {
		for o := range obstacles {
			grid[o.y][o.x] = 'O'
		}
		for _, line := range grid {
			fmt.Println(string(line))
		}
	}

	fmt.Println("part one:")
	fmt.Println(len(visited))
	fmt.Println("part two:")
	fmt.Println(len(obstacles))

}

func getData(f string) (Grid, State) {
	lines, _ := file.GetLines(f)
	grid := make(Grid, len(lines), len(lines))
	var y, x int

	found := false
	for j := 0; j < len(lines); j++ {
		grid[j] = []byte(lines[j])
		if !found {
			for i := 0; i < len(grid[j]); i++ {
				if grid[j][i] == '^' {
					found = true
					y = j
					x = i
					break
				}
			}
		}
	}

	return grid, State{x: x, y: y, d: NORTH}
}
