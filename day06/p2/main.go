package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/mbordner/aoc2024/common/geom"
)

// 976 too low
// 1147 too low
func main() {
	visitedCount := 1
	grid, x, y := getData("../test.txt")

	dir := geom.North
	for {

		if grid[y][x] == '.' {
			visitedCount++
		}
		if grid[y][x] != '^' && grid[y][x] != '+' {
			if dir == geom.North || dir == geom.South {
				if grid[y][x] == '-' {
					grid[y][x] = '+'
				} else {
					grid[y][x] = '|'
				}
			} else {
				if grid[y][x] == '|' {
					grid[y][x] = '+'
				} else {
					grid[y][x] = '-'
				}
			}
		}

		nx := x
		ny := y

		switch dir {
		case geom.North:
			ny--
		case geom.East:
			nx++
		case geom.South:
			ny++
		case geom.West:
			nx--
		}

		if ny < 0 || ny == len(grid) || nx < 0 || nx == len(grid[y]) {
			break
		} else if grid[ny][nx] == '#' {
			switch dir {
			case geom.North:
				dir = geom.East
			case geom.East:
				dir = geom.South
			case geom.South:
				dir = geom.West
			case geom.West:
				dir = geom.North
			}
		} else {
			y = ny
			x = nx
		}
	}

	for _, line := range grid {
		fmt.Println(string(line))
	}

	visited := make(map[geom.Pos[int]]geom.Pos[int])

	for y, gridLine := range grid {
		for x, b := range gridLine {
			if b == '+' || b == '^' || b == '|' || b == '-' {
				visited[geom.Pos[int]{X: x, Y: y, Z: 0}] = geom.Pos[int]{X: x, Y: y, Z: 0}
			}
		}
	}

	tls := make(geom.Positions[int], 0, len(visited))
	trs := make(geom.Positions[int], 0, len(visited))
	brs := make(geom.Positions[int], 0, len(visited))
	bls := make(geom.Positions[int], 0, len(visited))

	for p := range visited {
		x := p.X
		y := p.Y
		if y > 0 && grid[y-1][x] == '#' {
			tls = append(tls, p)
		}
		if x < len(grid[y])-1 && grid[y][x+1] == '#' {
			trs = append(trs, p)
		}
		if y < len(grid)-1 && grid[y+1][x] == '#' {
			brs = append(brs, p)
		}
		if x > 0 && grid[y][x-1] == '#' {
			bls = append(bls, p)
		}
	}

	// possibilities tl tr br  to find bl
	pTlTrBr := make([]geom.Positions[int], 0, len(visited))
	for _, tl := range tls {
		for _, tr := range trs {
			if tl.Y == tr.Y {
				for _, br := range brs {
					if tr.X == br.X {
						pTlTrBr = append(pTlTrBr, geom.Positions[int]{tl, tr, br})
					}
				}
			}
		}
	}

	// possibilities tr br bl  to find tl
	pTrBrBl := make([]geom.Positions[int], 0, len(visited))
	for _, tr := range trs {
		for _, br := range brs {
			if tr.X == br.X {
				for _, bl := range bls {
					if br.Y == bl.Y {
						pTrBrBl = append(pTrBrBl, geom.Positions[int]{tr, br, bl})
					}
				}
			}
		}
	}

	// possibilities br bl tl  to find tr
	pBrBlTl := make([]geom.Positions[int], 0, len(visited))
	for _, br := range brs {
		for _, bl := range bls {
			if br.Y == bl.Y {
				for _, tl := range tls {
					if bl.X == tl.X {
						pBrBlTl = append(pBrBlTl, geom.Positions[int]{br, bl, tl})
					}
				}
			}
		}
	}

	// possibilities bl tl tr  to find br
	pBlTlTr := make([]geom.Positions[int], 0, len(visited))
	for _, bl := range bls {
		for _, tl := range tls {
			if bl.X == tl.X {
				for _, tr := range trs {
					if tl.Y == tr.Y {
						pBlTlTr = append(pBlTlTr, geom.Positions[int]{bl, tl, tr})
					}
				}
			}
		}
	}

	fmt.Println("pTlTrBr", pTlTrBr)
	fmt.Println("pTrBrBl", pTrBrBl)
	fmt.Println("pBrBlTl", pBrBlTl)
	fmt.Println("pBlTlTr", pBlTlTr)

	fmt.Println("visitedCount", visitedCount)

	fmt.Println("searching for obstacle locations...")
	obstacles := 0
	for _, p := range pTlTrBr { // looking for bl
		if c, e := visited[geom.Pos[int]{X: p[0].X, Y: p[2].Y, Z: 0}]; e {
			if c.X > 0 {
				obstacles++
				fmt.Println(geom.Pos[int]{X: c.X - 1, Y: c.Y, Z: c.Z})
				grid[c.Y][c.X-1] = 'O'
			}
		}
	}

	for _, p := range pTrBrBl { // searching for tl
		if c, e := visited[geom.Pos[int]{X: p[2].X, Y: p[0].Y, Z: 0}]; e {
			if c.Y > 0 {
				obstacles++
				fmt.Println(geom.Pos[int]{X: c.X, Y: c.Y - 1, Z: c.Z})
				grid[c.Y-1][c.X] = 'O'
			}
		}
	}

	for _, p := range pBrBlTl { // searching for tr
		if c, e := visited[geom.Pos[int]{X: p[0].X, Y: p[2].Y, Z: 0}]; e {
			if c.X < len(grid[c.Y])-1 {
				obstacles++
				fmt.Println(geom.Pos[int]{X: c.X + 1, Y: c.Y, Z: c.Z})
				grid[c.Y][c.X+1] = 'O'
			}
		}
	}

	for _, p := range pBlTlTr { // searching for br
		if c, e := visited[geom.Pos[int]{X: p[2].X, Y: p[0].Y, Z: 0}]; e {
			if c.Y < len(grid)-1 {
				obstacles++
				fmt.Println(geom.Pos[int]{X: c.X, Y: c.Y + 1, Z: c.Z})
				grid[c.Y+1][c.X] = 'O'
			}

		}
	}

	for _, line := range grid {
		fmt.Println(string(line))
	}

	fmt.Println("total obstacles to place: ", obstacles)

}

func getData(f string) ([][]byte, int, int) {
	lines, _ := file.GetLines(f)
	grid := make([][]byte, len(lines), len(lines))
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

	return grid, x, y
}
