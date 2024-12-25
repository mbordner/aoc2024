package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strings"
)

type Object struct {
	bs [][]byte
	hs []int
}

func main() {
	keys, locks := getData("../data.txt")

	count := 0
	for _, k := range keys {
	nextLock:
		for _, l := range locks {
			for h := range k.hs {
				if k.hs[h]+l.hs[h] > 5 {
					continue nextLock
				}
			}
			count++
		}
	}

	fmt.Println(count)
}

func getData(f string) ([]*Object, []*Object) {
	content, _ := file.GetContent(f)

	objects := strings.Split(string(content), "\n\n")

	var keys []*Object  // bottom row #
	var locks []*Object // top row #

	for _, o := range objects {
		lines := strings.Split(o, "\n")
		obj := &Object{}
		obj.bs = make([][]byte, 0, len(lines))
		for _, l := range lines {
			obj.bs = append(obj.bs, []byte(l))
		}

		obj.hs = make([]int, len(lines[0]))

		if obj.bs[0][0] == '#' {
			locks = append(locks, obj)
			for x := 0; x < len(lines[0]); x++ {
				for y := 0; y < len(lines); y++ {
					if lines[y][x] != '#' {
						obj.hs[x] = y - 1
						break
					}
				}
			}
		} else {
			keys = append(keys, obj)
			for x := 0; x < len(lines[0]); x++ {
				for y := 0; y < len(lines); y++ {
					if lines[y][x] != '.' {
						obj.hs[x] = len(lines) - y - 1
						break
					}
				}
			}
		}
	}

	return keys, locks
}
