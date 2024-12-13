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

func main() {
	_, regions := getData("../data.txt")
	sum := 0
	for r := range regions {
		sum += r.area() * r.perimeter()
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
					n := &Plant{
						pos:       pos,
						region:    nil,
						species:   plant.species,
						neighbors: Plants{},
					}
					plant.neighbors = append(plant.neighbors, n)
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
