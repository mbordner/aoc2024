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
)

type Stones []string

func blink(stones Stones) Stones {
	next := make(Stones, 0, len(stones))

	for _, s := range stones {
		if reZeros.MatchString(s) {
			next = append(next, "1")
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
			next = append(next, Stones{l, r}...)
		} else {
			next = append(next, mul2024(s))
		}
	}

	return next
}

func mul2024(v string) string {
	n := new(big.Int)
	n, _ = n.SetString(v, 10)
	return n.Mul(n, big.NewInt(int64(2024))).String()
}

func main() {
	stones := getData("../data.txt")

	for i := 0; i < 25; i++ {
		stones = blink(stones)
	}

	fmt.Println(len(stones))

}

func getData(f string) Stones {
	data, _ := file.GetContent(f)
	return strings.Split(string(data), " ")
}
