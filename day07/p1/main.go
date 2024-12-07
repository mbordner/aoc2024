package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/bigexpression"
	"github.com/mbordner/aoc2024/common/file"
	"math/big"
	"regexp"
	"strings"
)

var (
	reEq = regexp.MustCompile(`(\d+):\s+(.*)`)
)

type Equations []Equation

type Equation struct {
	equation string
	value    *big.Int
}

func main() {
	eqs := getData("../data.txt")
	vars := make(map[string]*big.Int)

	precedence := func(op1, op2 string) int { return 0 }

	sum := big.NewInt(int64(0))

nextEQ:
	for _, eq := range eqs {
		equations := getEquations("+*", eq.equation)
		for _, equation := range equations {
			b, _ := bigexpression.NewParserWithPrecedence(equation, precedence)
			val := b.Eval(vars)
			if val.String() == eq.value.String() {
				sum = sum.Add(sum, val)
				continue nextEQ
			}
		}
	}

	fmt.Println(sum.String())
}

func getOpsRec(m map[string]bool, ops string, prefix string, l int) {
	if l == 0 {
		m[prefix] = true
		return
	}

	for i := 0; i < len(ops); i++ {
		nprefix := prefix + string(ops[i])
		getOpsRec(m, ops, nprefix, l-1)
	}
}

func getOps(ops string, l int) [][]string {
	m := make(map[string]bool)
	getOpsRec(m, ops, "", l)
	o := make([][]string, 0, len(m))
	for s := range m {
		o = append(o, strings.Split(s, ""))
	}
	return o
}

func getEquations(ops string, equation string) []string {
	tokens := strings.Split(equation, " ")
	l := len(tokens) - 1
	operators := getOps(ops, l)

	equations := make([]string, len(operators))
	for i := range operators {
		eq := tokens[0]
		for j := 0; j < len(tokens)-1; j++ {
			eq += operators[i][j] + tokens[j+1]
		}
		equations[i] = eq
	}

	return equations
}

func getData(f string) Equations {
	lines, _ := file.GetLines(f)
	eqs := make(Equations, len(lines))
	for i, line := range lines {
		matches := reEq.FindStringSubmatch(line)
		num := new(big.Int)
		num, _ = num.SetString(matches[1], 10)
		eqs[i] = Equation{equation: matches[2], value: num}
	}
	return eqs
}
