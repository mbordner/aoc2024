package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
)

type Plant struct {
	pos       common.Pos
	region    *Region
	species   string
	neighbors Plants
	fences    map[Dir]bool
}

type Dir int

const (
	UP Dir = iota
	RIGHT
	DOWN
	LEFT
)

type Line struct {
	p1 common.Pos
	p2 common.Pos
}

func (l *Line) String() string {
	return fmt.Sprintf("(%s->%s)", l.p1.String(), l.p2.String())
}

func (l *Line) dir() Dir {
	if l.p1.Y == l.p2.Y {
		if l.p1.X < l.p2.X {
			return RIGHT // top
		} else {
			return LEFT // bottom
		}
	} else {
		if l.p1.Y < l.p2.Y {
			return DOWN // right
		} else {
			return UP // left
		}
	}
}

func (l *Line) merge(o *Line) *Line {
	if o.dir() != l.dir() {
		return nil
	}
	s1, e1, s2, e2 := l.p1, l.p2, o.p1, o.p2
	var r Line
	switch l.dir() {
	case UP:
		if s1.X == s2.X {
			// if o fully contains l
			if s1.Y <= s2.Y && e1.Y >= e2.Y {
				r = *o
				return &r
			}
			// if l fully contains o
			if s2.Y <= s1.Y && e2.Y >= e1.Y {
				r = *l
				return &r
			}
			// if l start is on o
			// if o start is on l
			// if l end is on o
			// if o end is on l
			if (s1.Y <= s2.Y && s1.Y >= e2.Y) ||
				(s2.Y <= s1.Y && s2.Y >= e1.Y) ||
				(e1.Y <= s2.Y && e1.Y >= e2.Y) ||
				(e2.Y <= s1.Y && e2.Y >= e1.Y) {
				r.p1.X = s1.X
				r.p2.X = s1.X
				r.p1.Y = common.Max(s1.Y, s2.Y)
				r.p2.Y = common.Min(e1.Y, e2.Y)
				return &r
			}
		}
	case RIGHT:
		if s1.Y == s2.Y {
			// if o fully contains l
			if s1.X >= s2.X && e1.X <= e2.X {
				r = *o
				return &r
			}
			// if l fully contains o
			if s2.X >= s1.X && e2.X <= e1.X {
				r = *l
				return &r
			}
			// if l start is on o
			// if o start is on l
			// if l end is on o
			// if o end is on l
			if (s1.X >= s2.X && s1.X <= e2.X) ||
				(s2.X >= s1.X && s2.X <= e1.X) ||
				(e1.X >= s2.X && e1.X <= e2.X) ||
				(e2.X >= s1.X && e2.X <= e1.X) {
				r.p1.Y = s1.Y
				r.p2.Y = s1.Y
				r.p1.X = common.Min(s1.X, s2.X)
				r.p2.X = common.Max(e1.X, e2.X)
				return &r
			}

		}
	case DOWN:
		if s1.X == s2.X {
			// if o fully contains l
			if s1.Y >= s2.Y && e1.Y <= e2.Y {
				r = *o
				return &r
			}
			// if l fully contains o
			if s2.Y >= s1.Y && e2.Y <= e1.Y {
				r = *l
				return &r
			}
			// if l start is on o
			// if o start is on l
			// if l end is on o
			// if o end is on l
			if (s1.Y >= s2.Y && s1.Y <= e2.Y) ||
				(s2.Y >= s1.Y && s2.Y <= e1.Y) ||
				(e1.Y >= s2.Y && e1.Y <= e2.Y) ||
				(e2.Y >= s1.Y && e2.Y <= e1.Y) {
				r.p1.X = s1.X
				r.p2.X = s1.X
				r.p1.Y = common.Min(s1.Y, s2.Y)
				r.p2.Y = common.Max(e1.Y, e2.Y)
				return &r
			}
		}
	case LEFT:
		if s1.Y == s2.Y {
			// if o fully contains l
			if s1.X <= s2.X && e1.X >= e2.X {
				r = *o
				return &r
			}
			// if l fully contains o
			if s2.X <= s1.X && e2.X >= e1.X {
				r = *l
				return &r
			}
			// if l start is on o
			// if o start is on l
			// if l end is on o
			// if o end is on l
			if (s1.X <= s2.X && s1.X >= e2.X) ||
				(s2.X <= s1.X && s2.X >= e1.X) ||
				(e1.X <= s2.X && e1.X >= e2.X) ||
				(e2.X <= s1.X && e2.X >= e1.X) {
				r.p1.Y = s1.Y
				r.p2.Y = s1.Y
				r.p1.X = common.Max(s1.X, s2.X)
				r.p2.X = common.Min(e1.X, e2.X)
				return &r
			}
		}
	}
	return nil
}

type Garden map[common.Pos]*Plant
type Regions map[*Region]Plants

type Plants []*Plant

type Region struct {
	plants Plants
}

func (r *Region) area() int {
	return len(r.plants)
}

func (r *Region) perimeter() int {
	perimeter := 0
	for _, p := range r.plants {
		perimeter += 4 - len(p.neighbors)
	}
	return perimeter
}

func (r *Region) fences() []Line {
	starts := make(map[common.Pos]Line)
	for _, p := range r.plants {

		for d, e := range p.fences {
			if e {
				var pl Line
				switch d {
				case UP:
					// top: need RIGHT dir fence
					pl.p1.X = p.pos.X
					pl.p1.Y = p.pos.Y
					pl.p2.X = p.pos.X + 1
					pl.p2.Y = p.pos.Y
				case RIGHT:
					// need a DOWN dir fence
					pl.p1.X = p.pos.X + 1
					pl.p1.Y = p.pos.Y
					pl.p2.X = p.pos.X + 1
					pl.p2.Y = p.pos.Y + 1
				case DOWN:
					// bottom: need a LEFT dir fence
					pl.p1.X = p.pos.X + 1
					pl.p1.Y = p.pos.Y + 1
					pl.p2.X = p.pos.X
					pl.p2.Y = p.pos.Y + 1
				case LEFT:
					// need a UP dir fence
					pl.p1.X = p.pos.X
					pl.p1.Y = p.pos.Y + 1
					pl.p2.X = p.pos.X
					pl.p2.Y = p.pos.Y
				}
				starts[pl.p1] = pl
			}
		}

	}

	for {
		merged := false

		for start, l := range starts {
			if o, e := starts[l.p2]; e {
				r := (&l).merge(&o)
				if r != nil {
					starts[start] = *r
					delete(starts, l.p2)
					merged = true
					break
				}
			}
		}

		if !merged {
			break
		}
	}

	fences := make([]Line, 0, len(starts))

	for _, l := range starts {
		fences = append(fences, l)
	}

	/*
		sort.Slice(fences, func(i, j int) bool {
			if fences[i].p1.Y < fences[j].p1.Y {
				return true
			}
			if fences[i].p1.X < fences[j].p1.X {
				return true
			}
			return false
		})

		fmt.Println(fences)
	*/

	return fences
}

// 839354 too low
func main() {
	_, regions := getData("../data.txt")
	sum := 0
	for r := range regions {
		fences := r.fences()
		sum += r.area() * len(fences) // r.perimeter()
	}
	fmt.Println(sum)
}

func expand(grid common.Grid, garden Garden, plant *Plant) Plants {
	plants := Plants{plant}
	garden[plant.pos] = plant

	for y := plant.pos.Y - 1; y <= plant.pos.Y+1; y++ {
		for x := plant.pos.X - 1; x <= plant.pos.X+1; x++ {
			if (y != plant.pos.Y && x == plant.pos.X) || (x != plant.pos.X && y == plant.pos.Y) {
				pos := common.Pos{Y: y, X: x}
				if grid.ContainsPos(pos) && string(grid[y][x]) == plant.species {
					var d Dir
					if y == plant.pos.Y {
						if x < plant.pos.X {
							d = LEFT
						} else {
							d = RIGHT
						}
					} else {
						if y < plant.pos.Y {
							d = UP
						} else {
							d = DOWN
						}
					}
					n := &Plant{
						pos:       pos,
						region:    nil,
						species:   plant.species,
						neighbors: Plants{},
						fences:    map[Dir]bool{UP: true, RIGHT: true, DOWN: true, LEFT: true},
					}
					plant.neighbors = append(plant.neighbors, n)
					plant.fences[d] = false
					if _, e := garden[pos]; !e {
						plants = append(plants, expand(grid, garden, n)...)
					}
				}
			}
		}
	}

	return plants
}

func getData(f string) (Garden, Regions) {
	lines, _ := file.GetLines(f)
	grid := common.ConvertGrid(lines)

	garden := make(Garden)
	regions := make(Regions)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			pos := common.Pos{Y: y, X: x}

			if p, e := garden[pos]; !e {

				p = &Plant{
					pos:       pos,
					region:    nil,
					species:   string(grid[y][x]),
					neighbors: Plants{},
					fences:    map[Dir]bool{UP: true, RIGHT: true, DOWN: true, LEFT: true},
				}

				r := &Region{}
				r.plants = expand(grid, garden, p)

				for i := range r.plants {
					r.plants[i].region = r
				}

				regions[r] = r.plants

			}

		}
	}

	return garden, regions
}
