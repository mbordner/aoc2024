package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"math/big"
	"regexp"
	"strings"
)

var (
	reZeros    = regexp.MustCompile(`^0+$`)
	reUnpadded = regexp.MustCompile(`([1-9]\d*)`)
	counts     = Counts{}
)

type Stones []string

type Mem struct {
	s string
	n int
}
type Counts map[Mem]uint64

func (c Counts) add(m Mem, v uint64) {
	c[m] = v
}

func (c Counts) val(m Mem) uint64 {
	if v, e := c[m]; e {
		return v
	}
	return 0
}

func blink(s string, n int) uint64 {

	mem := Mem{s: s, n: n}
	val := counts.val(mem)

	if val > 0 {
		return val
	}

	if n == 0 {
		return 1
	}

	if reZeros.MatchString(s) {

		val = blink("1", n-1)

	} else if len(s)%2 == 0 {
		h := len(s) / 2
		l := s[0:h]
		r := s[h:]
		if reZeros.MatchString(r) {
			r = "0"
		} else {
			matches := reUnpadded.FindStringSubmatch(r)
			r = matches[1]
		}

		val = blink(l, n-1) + blink(r, n-1)

	} else {

		val = blink(mul2024(s), n-1)
	}

	counts[mem] = val
	return val
}

func mul2024(v string) string {
	n := new(big.Int)
	n, _ = n.SetString(v, 10)
	return n.Mul(n, big.NewInt(int64(2024))).String()
}

func main() {
	stones := getData("../data.txt")

	sum := uint64(0)

	for i := 0; i < len(stones); i++ {
		sum += blink(stones[i], 75)
	}

	fmt.Println(sum)

}

func getData(f string) Stones {
	data, _ := file.GetContent(f)
	return strings.Split(string(data), " ")
}
