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
	return s % 16777216 // takes last 24 bits
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

func test(s uint64) {
	fmt.Printf("%d: %064b\n", 0, s)
	fmt.Printf("%d: %064b shift left 6 bits\n", 1, s*64)
	n := mix(s, s*64)
	fmt.Printf("%d: %064b XOR #0 with first #1\n", 1, n)
	n = prune(n)
	fmt.Printf("%d: %064b get last 24 bits\n", 2, n)
	fmt.Printf("%d: %064b shift right 5 bits\n", 3, n/32)
	n = mix(n, n/32)
	fmt.Printf("%d: %064b XOR first #3 with #2\n", 3, n)
	n = prune(n)
	fmt.Printf("%d: %064b get last 24 bits\n", 4, n)
	fmt.Printf("%d: %064b shift left 11 bits\n", 5, n*2048)
	n = mix(n, n*2048)
	fmt.Printf("%d: %064b XOR first #5 with #4\n", 5, n)
	n = prune(n)
	fmt.Printf("%d: %064b get last 24 bits\n", 6, n)
}

type Seq struct {
	a int
	b int
	c int
	d int
}

type SeqTotals map[Seq]int

func (st SeqTotals) Has(s Seq) bool {
	if _, e := st[s]; e {
		return true
	}
	return false
}

func (st SeqTotals) Add(s Seq, v int) int {
	if st.Has(s) {
		st[s] += v
	} else {
		st[s] = v
	}
	return st[s]
}

func main() {

	data := getData("../data.txt")

	seqTotal := make(SeqTotals)
	maxPrice := 0
	var maxSeq Seq

	for _, sn := range data {
		buyerSeqTotals := make(SeqTotals)
		prices := make([]int, 2001)
		prices[0] = int(sn) % 10
		for p := 1; p < len(prices); p++ {
			sn = next(sn)
			prices[p] = int(sn) % 10
		}
		for i := 1; i < len(prices)-4; i++ {
			p := prices[i+3]
			s := Seq{}
			s.d, s.c, s.b, s.a = prices[i+3]-prices[i+2], prices[i+2]-prices[i+1], prices[i+1]-prices[i], prices[i]-prices[i-1]
			if !buyerSeqTotals.Has(s) {
				buyerSeqTotals.Add(s, p)
			}
		}
		for seq, price := range buyerSeqTotals {
			st := seqTotal.Add(seq, price)
			if st > maxPrice {
				maxPrice = st
				maxSeq = seq
			}
		}
	}

	fmt.Println(maxSeq)
	fmt.Println("max bananas: ", maxPrice)

}

func getData(f string) []uint64 {
	lines, _ := file.GetLines(f)
	data := make([]uint64, len(lines))
	for i, line := range lines {
		data[i], _ = strconv.ParseUint(line, 10, 64)
	}
	return data
}
