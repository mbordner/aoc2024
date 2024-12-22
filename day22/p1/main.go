package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"strconv"
)

func mix(s, n uint64) uint64 {
	return s ^ n
}

func prune(s uint64) uint64 {
	return s % 16777216
}

func next(s uint64) uint64 {
	n := mix(s, s*64)
	n = prune(n)
	n = mix(n, n/32)
	n = prune(n)
	n = mix(n, n*2048)
	n = prune(n)
	return n
}

func main() {
	//test1 := uint64(123)
	//for i := 0; i < 10; i++ {
	//	test1 = next(test1)
	//	fmt.Println(test1)
	//}

	data := getData("../data.txt")
	values := make([]uint64, len(data))
	copy(values, data)

	for i := 0; i < 2000; i++ {
		for j := range values {
			values[j] = next(values[j])
		}
	}

	sum := uint64(0)
	for _, v := range values {
		sum += v
	}

	fmt.Println(sum)
}

func getData(f string) []uint64 {
	lines, _ := file.GetLines(f)
	data := make([]uint64, len(lines))
	for i, line := range lines {
		data[i], _ = strconv.ParseUint(line, 10, 64)
	}
	return data
}
