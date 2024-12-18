package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"github.com/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	reProgram = regexp.MustCompile(`^Register A:\s+(\d+)\nRegister B:\s+(\d+)\nRegister C:\s+(\d+)\n\nProgram:\s+(.*)$`)
)

type Instruction int

// combos
// 0-3 are literal
// 4  val in RegA
// 5  val in RegB
// 6  val in RegC
// 7  reserved
const (
	adv Instruction = iota // div:  numerator in RegA, denominator is 2^(instruction combo op), result truncated, stored in RegA
	bxl                    // bitwise XOR RegB with (instruction literal operator)
	bst                    // (instruction combo op) mod 8
	jnz                    // does nothing if RegA is 0, otherwise advances ptr to (instruction literal operator)
	bxc                    // bitwise XOR of val in RegB and val in RegC
	out                    // (instruction combo op) mod 8 then outputs value
	bdv                    // like adv, but value stored in RegB
	cdv                    // like adv, but value stored in RegC
)

type Computer struct {
	a       int
	b       int
	c       int
	ptr     int
	program []int
	output  []int
}

func (c *Computer) Reset() {
	c.a, c.b, c.c, c.ptr = 0, 0, 0, 0
	c.output = make([]int, 0, 10)
}

func (c *Computer) Run() string {
	c.ptr = 0
	c.loop()
	return c.String()
}

func (c *Computer) String() string {
	output := make([]string, len(c.output))
	for i := range c.output {
		output[i] = strconv.Itoa(c.output[i])
	}
	return strings.Join(output, ",")
}

func (c *Computer) out(val int) {
	c.output = append(c.output, val)
}

func (c *Computer) xor(val1, val2 int) int {
	return val1 ^ val2
}

func (c *Computer) pow2(val int) int {
	return 1 << val
}

func (c *Computer) combo(val int) int {
	if val >= 0 && val <= 3 {
		return val
	}
	if val == 7 { // this is reserved, supposed to not happen
		panic(errors.New("invalid combo"))
	}
	switch val {
	case 4:
		return c.a
	case 5:
		return c.b
	case 6:
		return c.c
	}
	return 0
}

func (c *Computer) loop() {
nextInstr:
	for c.ptr < len(c.program) {
		op := Instruction(c.program[c.ptr])
		c.ptr++
		operand := c.program[c.ptr]
		c.ptr++
		//fmt.Println(c.ptr, fmt.Sprintf("[%d, %d]", op, operand), fmt.Sprintf("{A:%d,B:%d,C:%d}", c.a, c.b, c.c), c.String())

		regDiv := func(operand int) int {
			num := c.a
			div := c.pow2(c.combo(operand))
			return num / div
		}

		switch op {
		case adv:
			c.a = regDiv(operand)
		case bxl:
			c.b = c.xor(c.b, operand)
		case bst:
			c.b = c.combo(operand) % 8
		case jnz:
			if c.a != 0 {
				c.ptr = operand
				continue nextInstr
			}
		case bxc:
			c.b = c.xor(c.b, c.c)
		case out:
			c.out(c.combo(operand) % 8)
			match := 0
			for i := 0; i < len(c.output) && i < len(c.program); i++ {
				if c.output[i] == c.program[i] {
					match++
				}
			}
			if match != len(c.output) {
				break
			}
		case bdv:
			c.b = regDiv(operand)
		case cdv:
			c.c = regDiv(operand)
		}

	}
}

func exp(base, e int) int {
	return int(math.Pow(float64(base), float64(e)))
}

func main() {
	computer := getComputer("../data.txt")

	matched := 0
	a := 0
	for {
		computer.Reset()
		a++

		aStart := a
		aStart = a*exp(8, 9) + 0105714775
		computer.a = aStart

		computer.Run()
		match := 0
		for i := 0; i < len(computer.output) && i < len(computer.program); i++ {
			if computer.output[i] == computer.program[i] {
				match++
			}
		}
		if match == len(computer.program) {
			fmt.Println(aStart)
			break
		} else if match > matched {
			matched = match
			fmt.Println(aStart, fmt.Sprintf("%o", aStart), matched, computer.output)
		}
	}
	output := computer.Run()
	fmt.Println(output)
}

func getComputer(f string) *Computer {
	lines, _ := file.GetContent(f)
	matches := reProgram.FindStringSubmatch(string(lines))
	computer := &Computer{}
	computer.a = getIntVal(matches[1])
	computer.b = getIntVal(matches[2])
	computer.c = getIntVal(matches[3])
	computer.output = make([]int, 0, 10)
	tokens := strings.Split(matches[4], ",")
	computer.program = make([]int, len(tokens))
	for i, token := range tokens {
		computer.program[i] = getIntVal(token)
	}
	return computer
}

func getIntVal(s string) int {
	val, _ := strconv.ParseInt(s, 10, 32)
	return int(val)
}
