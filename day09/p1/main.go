package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
)

const (
	free = int64(-1)
)

type Disk []int64

func main() {
	data, _ := file.GetContent("../data.txt")
	id := int64(0)
	disk := make(Disk, 0, 1000)

	for i := 0; i < len(data); i++ {
		val := getIntVal(data[i])
		if i%2 == 0 {
			// file
			for j := int64(0); j < val; j++ {
				disk = append(disk, id)
			}
			id++
		} else {
			// free space
			for j := int64(0); j < val; j++ {
				disk = append(disk, free)
			}
		}
	}

	for i, j := 0, len(disk)-1; i < j; {
		for disk[i] != free && i < j {
			i++
		}
		for disk[j] == free && j > i {
			j--
		}
		for disk[i] == free && disk[j] != free && i < j {
			disk[i], disk[j] = disk[j], disk[i]
			i, j = i+1, j-1
		}
	}

	sum := int64(0)
	for i := 0; i < len(disk) && disk[i] != free; i++ {
		sum += int64(i) * disk[i]
	}

	fmt.Println(sum)
}

func getIntVal(c byte) int64 {
	val, _ := strconv.ParseInt(string(c), 10, 64)
	return val
}
