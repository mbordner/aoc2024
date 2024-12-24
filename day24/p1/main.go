package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type GateType int

func (t GateType) String() string {
	switch t {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	}
	return "UNKNOWN"
}

var (
	reWire = regexp.MustCompile(`(.*):\s+(1|0)\s*`)
	reGate = regexp.MustCompile(`(.*)\s+(AND|XOR|OR)\s+(.*)\s+->\s(.*)`)
	reZVal = regexp.MustCompile(`z0*(\d+)`)
)

const (
	UNKNOWN GateType = iota
	AND
	XOR
	OR
)

type WireValues map[string]int

type Wires map[string]chan int

type Gate struct {
	gateType   GateType
	inputNames []string
	outputName string
	done       chan<- bool
	inputWires []<-chan int
	outputWire chan<- int
}

func (g *Gate) String() string {
	return fmt.Sprintf("%s %s %s", g.inputNames[0], g.gateType.String(), g.inputNames[1])
}

func (g *Gate) Eval(wireValues WireValues) {
	var i0, i1 int

	for i := 0; i < 2; i++ {
		select {
		case i0 = <-g.inputWires[0]:

			fmt.Printf("gate (%s) received input 1: %d\n", g.String(), i0)
		case i1 = <-g.inputWires[1]:
			fmt.Printf("gate (%s) received input 2: %d\n", g.String(), i1)
		}
	}

	var val int

	switch g.gateType {
	case AND:
		val = i0 & i1
	case XOR:
		val = i0 ^ i1
	case OR:
		val = i0 | i1
	}

	wireValues[g.outputName] = val

	fmt.Printf("gate (%s) sent output (%s) val: %d\n", g.String(), g.outputName, val)

	g.outputWire <- val
	g.done <- true
}

func main() {
	initialWires, wires, gates, doneChannel := getData("../test.txt")

	var z int
	var wg sync.WaitGroup

	var zWires []string
	for w := range wires {
		if strings.HasPrefix(w, "z") {
			wg.Add(1)
			zWires = append(zWires, w)
		}
	}

	for _, w := range zWires {
		go func(wg *sync.WaitGroup, zChanName string, zChan <-chan int) {
			zBitVal := <-zChan
			matches := reZVal.FindStringSubmatch(zChanName)
			zBit := getIntVal(matches[1])
			zBitVal <<= zBit
			z |= zBitVal
			wg.Done()
		}(&wg, w, wires[w])
	}

	for n, v := range initialWires {
		wires[n] <- v
	}

	wg.Wait()

	for i := 0; i < len(gates); i++ {
		<-doneChannel
		fmt.Println("gate done")
	}

	fmt.Println("done")

	fmt.Println(z)
}

func getData(f string) (WireValues, Wires, []*Gate, chan bool) {

	wireValues := make(WireValues)

	doneChannel := make(chan bool)

	content, _ := file.GetContent(f)
	sections := strings.Split(string(content), "\n\n")
	for _, line := range strings.Split(sections[0], "\n") {
		if reWire.MatchString(line) {
			matches := reWire.FindStringSubmatch(line)
			wireValues[matches[1]] = getIntVal(matches[2])
		}
	}
	lines := strings.Split(sections[1], "\n")
	gates := make([]*Gate, len(lines))
	for i, line := range lines {
		if reGate.MatchString(line) {
			matches := reGate.FindStringSubmatch(line)
			gates[i] = &Gate{inputNames: []string{matches[1], matches[3]}, outputName: matches[4], done: doneChannel}
			switch matches[2] {
			case "AND":
				gates[i].gateType = AND
			case "XOR":
				gates[i].gateType = XOR
			case "OR":
				gates[i].gateType = OR
			}
		}

	}

	wires := make(Wires)

	for _, gate := range gates {
		gate.inputWires = make([]<-chan int, len(gate.inputNames))
		for i, input := range gate.inputNames {
			if _, e := wires[input]; !e {
				wires[input] = make(chan int)
			}
			gate.inputWires[i] = wires[input]
		}
		if _, e := wires[gate.outputName]; !e {
			wires[gate.outputName] = make(chan int)
		}
		gate.outputWire = wires[gate.outputName]

		go gate.Eval(wireValues)
	}

	return wireValues, wires, gates, doneChannel
}

func getIntVal(s string) int {
	val, _ := strconv.ParseInt(s, 10, 32)
	return int(val)
}
