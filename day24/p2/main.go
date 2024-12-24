package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common/file"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type GateType int

type WireValues map[string]int64

type Connector struct {
	wires  WireValues
	name   string
	inputs []chan int64
}

type Wires map[string]*Connector

type Gate struct {
	gateType        GateType
	inputNames      []string
	inputGates      Gates
	outputName      string
	done            chan<- bool
	inputWires      []<-chan int64
	outputConnector *Connector
}

type Gates []*Gate

var (
	debug    = false
	reWire   = regexp.MustCompile(`(.*):\s+(1|0)\s*`)
	reGate   = regexp.MustCompile(`(.*)\s+(AND|XOR|OR)\s+(.*)\s+->\s(.*)`)
	reBitPos = regexp.MustCompile(`^\D+0*(\d+)`)
	mu       sync.Mutex
)

const (
	UNKNOWN GateType = iota
	AND
	XOR
	OR
)

func (gs Gates) GetInputGate(outputName string) *Gate {
	for _, g := range gs {
		if g.outputName == outputName {
			return g
		}
	}
	return nil
}

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

func (c *Connector) Send(val int64) {
	mu.Lock()
	c.wires[c.name] = val
	for i := range c.inputs {
		c.inputs[i] <- val
	}
	mu.Unlock()
}

func (wv WireValues) SetBitMap(prefix string, val int64) {

	for bitName := range wv {
		if strings.HasPrefix(bitName, prefix) {
			matches := reBitPos.FindStringSubmatch(bitName)
			bitPos := getIntVal(matches[1])

			if (val & (1 << bitPos)) != 0 {
				wv[bitName] = int64(1)
			} else {
				wv[bitName] = int64(0)
			}
		}
	}

}

func (wv WireValues) GetBitMap(prefix string) int64 {

	type bit struct {
		pos int64
		val int64
	}
	var bits []bit

	for bitName, v := range wv {
		if strings.HasPrefix(bitName, prefix) {
			matches := reBitPos.FindStringSubmatch(bitName)
			b := bit{pos: getIntVal(matches[1]), val: v}
			bits = append(bits, b)
		}
	}

	sort.Slice(bits, func(i, j int) bool {
		return bits[i].pos > bits[j].pos
	})

	val := int64(0)

	for _, b := range bits {
		bval := b.val
		bval <<= b.pos
		val |= bval
	}

	return val
}

func (g *Gate) String() string {
	return fmt.Sprintf("%s %s %s -> %s", g.inputNames[0], g.gateType.String(), g.inputNames[1], g.outputConnector.name)
}

func (g *Gate) Eval(wireValues WireValues) {
	var i0, i1 int64

	for i := 0; i < 2; i++ {
		select {
		case i0 = <-g.inputWires[0]:
			if debug {
				fmt.Printf("gate (%s) received input 1: %d\n", g.String(), i0)
			}
		case i1 = <-g.inputWires[1]:
			if debug {
				fmt.Printf("gate (%s) received input 2: %d\n", g.String(), i1)
			}
		}
	}

	var val int64

	switch g.gateType {
	case AND:
		val = i0 & i1
	case XOR:
		val = i0 ^ i1
	case OR:
		val = i0 | i1
	}

	if debug {
		fmt.Printf("gate (%s) sent output (%s) val: %d\n", g.String(), g.outputName, val)
	}

	g.done <- true
	g.outputConnector.Send(val)
}

func calculate(wireValues WireValues, wires Wires, gates []*Gate, doneChannel chan bool) int64 {
	var z int64
	var wg1, wg2 sync.WaitGroup

	for _, gate := range gates {
		go gate.Eval(wireValues)
	}

	var zWires []string
	for w := range wires {
		if strings.HasPrefix(w, "z") {
			wg1.Add(1)
			wg2.Add(1)
			zWires = append(zWires, w)
		}
	}

	for i := range zWires {
		zWireName := zWires[i]
		zInput := make(chan int64)
		wires[zWireName].inputs = []chan int64{zInput}
		matches := reBitPos.FindStringSubmatch(zWireName)
		zBit := getIntVal(matches[1])
		go func(wg *sync.WaitGroup, wg2 *sync.WaitGroup, zWireName string, zBit int64, zChan chan int64) {
			wg1.Done()
			zBitVal := <-zChan
			if debug {
				fmt.Println("got z value:", zWireName, zBitVal)
			}
			zBitVal <<= zBit
			z |= zBitVal
			wg2.Done()
		}(&wg1, &wg2, zWireName, zBit, zInput)
	}

	wg1.Wait()

	for n, v := range wireValues {
		wires[n].Send(v)
	}

	for i := 0; i < len(gates); i++ {
		<-doneChannel
		if debug {
			fmt.Println("gate done")
		}
	}

	wg2.Wait()

	if debug {
		fmt.Println("done")
	}

	return z
}

func getData(f string) (WireValues, Wires, Gates, chan bool) {

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
	gates := make(Gates, len(lines))
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

	for wire := range wireValues {
		if _, e := wires[wire]; !e {
			wires[wire] = &Connector{wires: wireValues, name: wire, inputs: []chan int64{}}
		}
	}

	for _, gate := range gates {
		if _, e := wires[gate.outputName]; !e {
			wires[gate.outputName] = &Connector{wires: wireValues, name: gate.outputName, inputs: []chan int64{}}
			gate.outputConnector = wires[gate.outputName]
		}
	}

	for _, gate := range gates {
		gate.inputWires = make([]<-chan int64, len(gate.inputNames))
		gate.inputGates = make(Gates, len(gate.inputNames))
		for i, input := range gate.inputNames {
			intChan := make(chan int64)
			gate.inputWires[i] = intChan
			wires[input].inputs = append(wires[input].inputs, intChan)
			gate.inputGates[i] = gates.GetInputGate(input)
		}
	}

	return wireValues, wires, gates, doneChannel
}

func getIntVal(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func main() {
	data()
}

func data() {
	initialWires, wires, gates, doneChannel := getData("../data.txt")

	for g := range gates {
		fmt.Printf("gate %03d: %s\n", g, gates[g].String())
	}

	initialWires.SetBitMap("x", int64(70368744177663)) // 2^46 - 1
	initialWires.SetBitMap("y", int64(70368744177663))

	x := initialWires.GetBitMap("x")
	y := initialWires.GetBitMap("y")

	fmt.Printf("x: %064b\ny: %064b\n", x, y)

	zGates := make(Gates, 46)
	for i := range zGates {
		zName := fmt.Sprintf("z%02d", i)
		zGates[i] = gates.GetInputGate(zName)
	}

	z := calculate(initialWires, wires, gates, doneChannel)

	fmt.Printf("z: %064b\n", z)
	fmt.Printf("done")
}

func test3() {
	initialWires, wires, gates, doneChannel := getData("../test3.txt")

	gates[0].outputConnector, gates[5].outputConnector = gates[5].outputConnector, gates[0].outputConnector
	gates[1].outputConnector, gates[2].outputConnector = gates[2].outputConnector, gates[1].outputConnector

	x := initialWires.GetBitMap("x")
	y := initialWires.GetBitMap("y")

	z := calculate(initialWires, wires, gates, doneChannel)

	fmt.Printf("x: %064b\ny: %064b\nz: %064b\n", x, y, z)
}
