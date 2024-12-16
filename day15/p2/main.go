package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"io"
	"log"
	"strings"
	"time"
)

const (
	robotChar    = byte('@')
	wallChar     = byte('#')
	boxChar      = byte('O')
	openChar     = byte('.')
	upChar       = byte('^')
	rightChar    = byte('>')
	downChar     = byte('v')
	leftChar     = byte('<')
	unknownChar  = byte('?')
	leftBoxChar  = byte('[')
	rightBoxChar = byte(']')
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	Gray    = Color("\033[1;37m%s\033[0m")
	White   = Color("\033[1;97m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type Warehouse struct {
	grid      common.Grid
	robot     *Robot
	boxCells  map[common.Pos]int
	boxes     []common.Positions
	openCells map[common.Pos]bool
	ss        bool
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

func (wh *Warehouse) hasBox(p common.Pos) bool {
	if _, e := wh.boxCells[p]; e {
		return true
	}
	return false
}

func (wh *Warehouse) movePositions(id int, dir byte) map[int]common.Positions {
	ps := wh.boxes[id]
	nps := make(common.Positions, len(ps))
	for i := range ps {
		nps[i] = wh.getNextPos(ps[i], dir)
	}

	for _, np := range nps {
		if wh.openCells[np] == false {
			return nil
		}
	}

	toMove := make(map[int]common.Positions)
	for _, np := range nps {
		if wh.hasBox(np) {
			npId := wh.boxCells[np]
			if _, e := toMove[npId]; !e {
				npToMove := wh.movePositions(npId, dir)
				if npToMove == nil {
					return nil
				}
				for tid, tps := range npToMove {
					toMove[tid] = tps
				}
			}
		}
	}

	toMove[id] = nps

	return toMove
}

// p is position we want to move, dir is dir to move it
func (wh *Warehouse) movePos(p common.Pos, dir byte) *common.Pos {
	np := wh.getNextPos(p, dir)
	if wh.openCells[np] == false {
		return nil
	} else {
		if wh.hasBox(np) {

			if len(wh.boxes[wh.boxCells[p]]) > 1 {
				id := wh.boxCells[np]
				if dir == leftChar || dir == rightChar {
					headPos := wh.boxes[id][0]
					tailPos := wh.boxes[id][1]
					if dir == rightChar {
						headPos, tailPos = tailPos, headPos
					}

					nbp := wh.movePos(headPos, dir)
					if nbp == nil {
						return nil
					}

					for _, tp := range wh.boxes[id] {
						delete(wh.boxCells, tp)
					}

					if dir == leftChar {
						wh.boxes[id] = common.Positions{*nbp, headPos}
					} else {
						wh.boxes[id] = common.Positions{headPos, *nbp}
					}

					for _, tp := range wh.boxes[id] {
						wh.boxCells[tp] = id
					}

				} else {

					nps := wh.movePositions(id, dir)
					if nps == nil {
						return nil
					}

					for tid, tps := range nps {
						for _, tp := range wh.boxes[tid] {
							delete(wh.boxCells, tp)
						}
						wh.boxes[tid] = tps
					}

					for tid, _ := range nps {
						for _, tp := range wh.boxes[tid] {
							wh.boxCells[tp] = tid
						}
					}

				}
			} else {

				nbp := wh.movePos(np, dir)

				if nbp == nil {
					return nil // couldn't move the box at np
				}

				wh.boxCells[*nbp] = wh.boxCells[np]
				wh.boxes[wh.boxCells[np]][0] = *nbp
				delete(wh.boxCells, np)

			}

		}
	}
	return &np
}

func (wh *Warehouse) moveRobot() bool {
	nd := string(wh.robot.nextMove())
	if nd[0] != unknownChar {
		if np := wh.movePos(wh.robot.p, nd[0]); np != nil {
			wh.robot.commitMove(np)
		} else {
			wh.robot.commitMove(&wh.robot.p)
		}
		return true
	}
	return false
}

func getData(f string, superSize bool) *Warehouse {
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

	boxCells := make(map[common.Pos]int)
	nextBoxId := 0

	openCells := make(map[common.Pos]bool)

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			p := common.Pos{Y: y, X: x}
			if grid[y][x] == robotChar {
				robot.p.Y, robot.p.X = y, x
				grid[y][x] = openChar
			} else if grid[y][x] == boxChar {
				boxCells[p] = nextBoxId
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

	boxes := make([]common.Positions, len(boxCells))
	for p, i := range boxCells {
		boxes[i] = common.Positions{p}
	}

	if superSize {
		ssOpenCells := make(map[common.Pos]bool)
		ssGrid := make(common.Grid, len(grid))
		for y := range grid {
			ssGrid[y] = make([]byte, len(grid[y])*2)
			for x, i := 0, 0; x < len(grid[y]); x, i = x+1, i+2 {
				ps := common.Positions{common.Pos{Y: y, X: i}, common.Pos{Y: y, X: i + 1}}
				for _, p := range ps {
					if grid[y][x] == openChar {
						ssGrid[y][p.X] = openChar
						ssOpenCells[p] = true
					} else {
						ssGrid[y][p.X] = wallChar
						ssOpenCells[p] = false
					}
				}

			}
		}

		ssBoxCells := make(map[common.Pos]int)
		for p, i := range boxCells {
			boxes[i] = common.Positions{common.Pos{Y: p.Y, X: p.X * 2}, common.Pos{Y: p.Y, X: p.X*2 + 1}}
			ssBoxCells[boxes[i][0]] = i
			ssBoxCells[boxes[i][1]] = i

		}

		grid = ssGrid
		openCells = ssOpenCells
		boxCells = ssBoxCells
		robot.p.X *= 2
	}

	warehouse := Warehouse{grid: grid, robot: &robot, boxes: boxes, boxCells: boxCells, openCells: openCells, ss: superSize}
	return &warehouse
}

func (wh *Warehouse) print(out io.ReadWriter) {
	grid := make([][]string, len(wh.grid))
	for y := range wh.grid {
		grid[y] = make([]string, len(wh.grid[y]))
		for x := 0; x < len(wh.grid[y]); x++ {
			if wh.openCells[common.Pos{Y: y, X: x}] {
				grid[y][x] = Red(".")
			} else {
				grid[y][x] = Green("#")
			}
		}
	}

	if !wh.ss {
		for p := range wh.boxCells {
			grid[p.Y][p.X] = Yellow(string(boxChar))
		}
	} else {
		for _, ps := range wh.boxes {
			grid[ps[0].Y][ps[0].X] = Yellow(string(leftBoxChar))
			grid[ps[1].Y][ps[1].X] = Yellow(string(rightBoxChar))
		}
	}

	grid[wh.robot.p.Y][wh.robot.p.X] = White(string(robotChar)) // wh.robot.nextMove(
	for _, line := range grid {
		fmt.Fprintln(out, strings.Join(line, ""))
	}
}

func (wh *Warehouse) gpsSum() int {
	sum := 0
	for _, ps := range wh.boxes {
		if wh.ss {
			if len(ps) != 2 || ps[1].X != ps[0].X+1 {
				fmt.Println("something wrong?")
			}
		} else {
			if len(ps) != 1 {
				fmt.Println("something wrong?")
			}
		}
		//fmt.Println(ps[0], ps[0].Y*100+ps[0].X)
		sum += ps[0].Y*100 + ps[0].X
	}
	return sum
}

var stepThrough = false

// 1404329 too high
func main() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func loop(g *gocui.Gui) {
	wh := getData("../data.txt", true)

	for {

		moved := wh.moveRobot()
		if !moved {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("warehouse")
				fmt.Fprintln(v, "GPS Sum calculated from ", len(wh.boxes), "boxes:", wh.gpsSum())
				return nil
			})
			break
		}

		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("warehouse")
			v.Clear()
			wh.print(v)
			return nil
		})

		time.Sleep(time.Millisecond * 4)

	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX > 0 && maxY > 0 {
		if v, err := g.SetView("warehouse", 0, 0, 150, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			go loop(g)
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
