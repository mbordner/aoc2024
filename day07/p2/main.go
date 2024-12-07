package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
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

var (
	numTotalEquations = 0
	numProcessed      = 0
	numFiltered       = 0
)

func main() {
	eqs := getData("../data.txt")
	vars := make(map[string]*big.Int)

	precedence := func(op1, op2 string) int { return 0 }

	sum := big.NewInt(int64(0))

nextEQ:
	for _, eq := range eqs {
		equations := getEquations("+*|", eq.equation)
		numTotalEquations += len(equations)
		for _, equation := range equations {
			if validate(eq, equation) {
				b, _ := bigexpression.NewParserWithPrecedence(equation, precedence)
				val := b.Eval(vars)
				numProcessed++
				if val.String() == eq.value.String() {
					sum = sum.Add(sum, val)
					continue nextEQ
				}
			} else {
				numFiltered++
			}
		}
	}

	fmt.Println("total possible equations:", numTotalEquations)
	fmt.Println("equations evaluated:", numProcessed)
	fmt.Println("equations filtered:", numFiltered)
	fmt.Println("sum:", sum.String())
	// 492383931650959

}

func validate(eq Equation, equation string) bool {
	splitter := func(r rune) bool {
		return r == '+' || r == '*'
	}
	tokens := strings.FieldsFunc(equation, splitter)
	for _, token := range tokens {
		token = strings.Join(strings.Split(token, "|"), "")
		n := new(big.Int)
		n, _ = n.SetString(token, 10)
		if n.Cmp(eq.value) > 0 {
			return false
		}
	}
	return true
}

func getOps(ops string, l int) [][]string {
	m := make(map[string]bool)
	common.PopulateStringCombinationsAtLength(m, ops, "", l)
	o := make([][]string, 0, len(m))
	for s := range m {
		o = append(o, strings.Split(s, ""))
	}
	return o
}

func getEquations(ops string, equation string) []string {
	tokens := strings.Split(equation, " ")
	l := len(tokens) - 1
	operators := getOps(ops, l) // returns array of string arrays
	// len(operators) == the number of unique equations that can be generated from the " " gaps that can become
	// 		an expression operator
	// each operators element will be an array of string operators being the operator to use at the gaps in order
	//
	// this length should be length, num of unique operators (or length of ops) raised to power of number of gaps
	// e.g. if ops "+*" is length 2, and we have 3 space characters in the equation (3 spots where we can add ops)
	// this will be 2^3 (8) unique equation combinations

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
