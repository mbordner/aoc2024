package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
)

const (
	robotChar   = byte('@')
	wallChar    = byte('#')
	boxChar     = byte('O')
	openChar    = byte('.')
	upChar      = byte('^')
	rightChar   = byte('>')
	downChar    = byte('v')
	leftChar    = byte('<')
	unknownChar = byte('?')
)

type Warehouse struct {
	grid      common.Grid
	robot     *Robot
	boxes     map[common.Pos]int
	openCells map[common.Pos]bool
}

type Robot struct {
	p     common.Pos
	moves []byte
	mPtr  int
}

func (r *Robot) nextMove() byte {
	move := unknownChar
	if r.mPtr < len(r.moves) {
		move = r.moves[r.mPtr]
	}
	return move
}

func (r *Robot) commitMove(p *common.Pos) {
	r.p.X, r.p.Y = p.X, p.Y
	r.mPtr++
}

func (wh *Warehouse) getNextPos(p common.Pos, dir byte) common.Pos {
	switch dir {
	case upChar:
		return common.Pos{X: p.X, Y: p.Y - 1}
	case rightChar:
		return common.Pos{X: p.X + 1, Y: p.Y}
	case downChar:
		return common.Pos{X: p.X, Y: p.Y + 1}
	case leftChar:
		return common.Pos{X: p.X - 1, Y: p.Y}
	}
	return p
}

func (wh *Warehouse) print() {
	grid := make(common.Grid, len(wh.grid))
	for y := range wh.grid {
		grid[y] = make([]byte, len(wh.grid[y]))
		copy(grid[y], wh.grid[y])
	}
	for p := range wh.boxes {
		grid[p.Y][p.X] = boxChar
	}
	grid[wh.robot.p.Y][wh.robot.p.X] = robotChar
	for _, line := range grid {
		fmt.Println(string(line))
	}
}

func (wh *Warehouse) hasBox(p common.Pos) bool {
	if _, e := wh.boxes[p]; e {
		return true
	}
	return false
}

// p is position we want to move, dir is dir to move it
func (wh *Warehouse) movePos(p common.Pos, dir byte) *common.Pos {
	np := wh.getNextPos(p, dir)
	if wh.openCells[np] == false {
		return nil
	} else {
		if wh.hasBox(np) {
			nbp := wh.movePos(np, dir)
			if nbp == nil {
				return nil // couldn't move the box at np
			}
			wh.boxes[*nbp] = wh.boxes[np]
			delete(wh.boxes, np)
		}
	}
	return &np
}

func (wh *Warehouse) moveRobot() bool {
	nd := wh.robot.nextMove()
	if nd != unknownChar {
		if np := wh.movePos(wh.robot.p, nd); np != nil {
			wh.robot.commitMove(np)
		} else {
			wh.robot.commitMove(&wh.robot.p)
		}
		return true
	}
	return false
}

func (wh *Warehouse) gpsSum() int {
	sum := 0
	for p := range wh.boxes {
		sum += p.Y*100 + p.X
	}
	return sum
}

func main() {
	wh := getData("../data.txt")

	//wh.print()
	//reader := bufio.NewReader(os.Stdin)

	for {
		//_, _, _ = reader.ReadRune()
		moved := wh.moveRobot()
		if !moved {
			break
		}
		//fmt.Println("-==-------=-=-=--------")
		//wh.print()
	}

	wh.print()
	fmt.Println(wh.gpsSum())

}

func getData(f string) *Warehouse {
	var grid common.Grid

	lines, _ := file.GetLines(f)
	for i := range lines {
		if lines[i] == "" {
			grid = make(common.Grid, i)
			for j := 0; j < i; j++ {
				grid[j] = []byte(lines[j])
			}
			lines = lines[i+1:]
			break
		}
	}

	size := 0
	for _, line := range lines {
		size += len(line)
	}

	moves := make([]byte, 0, size)
	for _, line := range lines {
		moves = append(moves, []byte(line)...)
	}

	robot := Robot{
		p:     common.Pos{},
		moves: moves,
		mPtr:  0,
	}

	boxes := make(map[common.Pos]int)
	nextBoxId := 0

	openCells := make(map[common.Pos]bool)

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			p := common.Pos{Y: y, X: x}
			if grid[y][x] == robotChar {
				robot.p.Y, robot.p.X = y, x
				grid[y][x] = openChar
			} else if grid[y][x] == boxChar {
				boxes[p] = nextBoxId
				nextBoxId++
				grid[y][x] = openChar
			}
			if grid[y][x] == openChar {
				openCells[p] = true
			} else {
				openCells[p] = false
			}
		}
	}

	warehouse := Warehouse{grid: grid, robot: &robot, boxes: boxes, openCells: openCells}
	return &warehouse
}
