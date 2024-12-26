package main

import (
	"fmt"
	"github.com/mbordner/aoc2024/common"
	"github.com/mbordner/aoc2024/common/file"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type GatePairs map[*Gate]*Gate

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

func (g *Gate) PrintWithWireValues(wireValues WireValues) {

	indent := func(c int) string {
		return strings.Repeat("  ", c)
	}

	var gateTreePrinter func(g *Gate, depth int)

	gateTreePrinter = func(g *Gate, depth int) {
		ia := depth * 2
		fmt.Printf("%s%s(%d) <- %s (%p) <- \n", indent(ia), g.outputName, wireValues[g.outputName], g.gateType.String(), g)
		i, j := 0, 1
		if g.inputGates[j] != nil {
			if (g.inputGates[j].gateType == XOR || g.inputGates[j].gateType == OR) && g.inputGates[i].gateType != XOR {
				j, i = i, j
			} else if g.inputGates[j].gateType == AND && g.inputGates[i].gateType == AND && (g.inputGates[j].inputGates[0] != nil || g.inputGates[j].inputGates[1] != nil) {
				j, i = i, j
			}

		}

		if g.inputGates[i] != nil {
			gateTreePrinter(g.inputGates[i], depth+1)
		} else {
			fmt.Printf("%s<- %s(%d)\n", indent((depth+1)*2+2), g.inputNames[i], wireValues[g.inputNames[i]])
		}
		if g.inputGates[j] != nil {
			gateTreePrinter(g.inputGates[j], depth+1)
		} else {
			fmt.Printf("%s<- %s(%d)\n", indent((depth+1)*2+2), g.inputNames[j], wireValues[g.inputNames[j]])
		}
	}

	gateTreePrinter(g, 0)
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

func getBitArray(val int64, bits int) []int64 {
	vals := make([]int64, bits)
	for i := 0; i < bits; i++ {
		if (val & (1 << i)) != 0 {
			vals[i] = 1
		} else {
			vals[i] = 0
		}
	}
	return vals
}

func getAllBitGates(root *Gate) Gates {
	gates := Gates{root}
	for _, ig := range root.inputGates {
		if ig != nil {
			gates = append(gates, getAllBitGates(ig)...)
		}
	}
	return gates
}

func solve1() {
	initialWires, wires, gates, doneChannel := getData("../data.txt")

	numBits := 46
	expected := int64(math.Pow(2, float64(numBits)))
	expectedBits := getBitArray(expected, numBits)

	initialWires.SetBitMap("x", expected-1) // 2^46 - 1
	initialWires.SetBitMap("y", expected-1)

	zGates := make(Gates, numBits)
	for i := range zGates {
		zName := fmt.Sprintf("z%02d", i)
		zGates[i] = gates.GetInputGate(zName)
	}

	swaps := make(GatePairs)

	z := int64(0)

	bPtr := 0

	for z != expected {
		z = calculate(initialWires, wires, gates, doneChannel)
		zs := getBitArray(z, 46)
		for ; bPtr < len(zs); bPtr++ {
			if zs[bPtr] != expectedBits[bPtr] {

				bitGate := zGates[bPtr]
				bitGates := getAllBitGates(bitGate)

				bigGatePairSets := common.GetPairSets(bitGates)
				for _, p := range bigGatePairSets {
					if slices.Contains(p[0].inputNames, p[1].outputConnector.name) ||
						slices.Contains(p[1].inputNames, p[0].outputConnector.name) {
						continue // would create a loop
					}
					p[0].outputConnector, p[1].outputConnector = p[1].outputConnector, p[0].outputConnector
					tz := calculate(initialWires, wires, gates, doneChannel)
					tzs := getBitArray(tz, 46)
					if tzs[bPtr] == expectedBits[bPtr] {
						swaps[p[0]] = p[1]
						z = tz
						zs = tzs
						break
					} else {
						// swap back
						p[1].outputConnector, p[0].outputConnector = p[0].outputConnector, p[1].outputConnector
					}
				}

				fmt.Println(len(bigGatePairSets))
			}
		}

	}

	fmt.Println(len(swaps))

	fmt.Printf("z: %064b\n", z)
	fmt.Printf("done")
}

func (gs Gates) swap(wires Wires, g1output, g2output string) {
	var g1, g2 *Gate
	for g := range gs {
		if gs[g].outputName == g1output {
			g1 = gs[g]
		}
		if gs[g].outputName == g2output {
			g2 = gs[g]
		}
	}

	g1.outputName, g2.outputName = g2.outputName, g1.outputName
	g1.outputConnector, g2.outputConnector = g2.outputConnector, g1.outputConnector

	for _, g := range gs {
		for inputIndex, inputName := range g.inputNames {
			g.inputGates[inputIndex] = gs.GetInputGate(inputName)
		}
	}
}

func data() {
	initialWires, wires, gates, doneChannel := getData("../data.txt")

	gates.swap(wires, "hmk", "z16")

	gates.swap(wires, "z20", "fhp")

	gates.swap(wires, "tpc", "rvf")

	gates.swap(wires, "z33", "fcd")

	//gates.swap(wires, "pgp", "scf")

	//gates.swap(wires, "fmg", "mbp")
	//gates.swap(wires, "z33", "smf")

	//gates.swap(wires, "hmk", "z16") // at test bit 16
	//gates.swap(wires, "hwg", "wbg")
	//gates.swap(wires, "qpp", "qgq")
	//gates.swap(wires, "fdg", "tqw")
	//gates.swap(wires, "mpw", "fdg")
	//gates.swap(wires, "tqw", "mpw")
	//gates.swap(wires, "tqv", "nbj")
	//gates.swap(wires, "hmk", "ndj")
	//gates.swap(wires, "fhp", "z20")
	//gates.swap(wires, "ncf", "vcj")
	//gates.swap(wires, "tpc", "rvf")

	//for g := range gates {
	//	fmt.Printf("gate %03d: %s\n", g, gates[g].String())
	//}

	testBit := 45
	initialWires.SetBitMap("y", int64(1)<<testBit-1) // 2^46 - 1
	initialWires.SetBitMap("x", int64(1)<<testBit-1)

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

	for i := range zGates {
		fmt.Printf("============== z%02d ====================\n", i)
		zGates[i].PrintWithWireValues(initialWires)
	}

}

func test3() {
	initialWires, wires, gates, doneChannel := getData("../test3.txt")

	gates.swap(wires, "z00", "z05")
	gates.swap(wires, "z01", "z02")

	x := initialWires.GetBitMap("x")
	y := initialWires.GetBitMap("y")

	zGates := make(Gates, 6)
	for i := range zGates {
		zName := fmt.Sprintf("z%02d", i)
		zGates[i] = gates.GetInputGate(zName)
	}

	z := calculate(initialWires, wires, gates, doneChannel)

	fmt.Printf("x: %064b\ny: %064b\nz: %064b\n", x, y, z)

	for i := range zGates {
		fmt.Printf("============== z%02d ====================\n", i)
		zGates[i].PrintWithWireValues(initialWires)
	}
}

func main() {
	data()
}
