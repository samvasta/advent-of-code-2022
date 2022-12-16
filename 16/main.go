package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func parseInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		panic("wtf")
	}

	return v
}

type System struct {
	connections map[string][]string

	valves map[string]int
}

type Minute struct {
	prev         *Minute
	minutesLeft  int
	move1, move2 *Move
}
type Move struct {
	prev *Move

	minutesLeft int
	openValve   bool
	valve       string

	valvePressure int
}

func (m *Move) totalPressureReleased() int {
	if m.openValve {
		return m.minutesLeft * m.valvePressure
	}
	return 0
}

func (m *Minute) totalPressureReleasedMin() int {
	if m.move2 != nil {
		return m.move1.totalPressureReleased() + m.move1.totalPressureReleased()
	}
	return m.move1.totalPressureReleased()
}

func (m *Move) cumulativePressureReleased() int {
	if m.prev == nil {
		return m.totalPressureReleased()
	}
	return m.totalPressureReleased() + m.prev.cumulativePressureReleased()
}
func (m *Minute) cumulativePressureReleasedMin() int {
	if m.prev == nil {
		return m.totalPressureReleasedMin()
	}
	return m.totalPressureReleasedMin() + m.prev.cumulativePressureReleasedMin()
}

func (m *Move) String() string {
	str := ""

	if m.prev != nil {
		str += m.prev.String()
	}
	action := "leave valve closed"
	if m.openValve {
		action = "open valve"
	}

	str += fmt.Sprintf("[%v] At %v, %v releasing %v (%v*%v) pressure (cumulative=%v)\n", m.minutesLeft, m.valve, action, m.totalPressureReleased(), m.valvePressure, m.minutesLeft, m.cumulativePressureReleased())

	return str
}
func (m *Move) Visited(v string) bool {
	if m.valve == v {
		return true
	}
	if m.prev == nil {
		return false
	}
	return m.prev.Visited(v)
}

func (m *Move) IsOpen(v string) bool {
	if m.valve == v && m.openValve {
		return true
	}
	if m.prev == nil {
		return false
	}
	return m.prev.IsOpen(v)
}
func (m *Minute) VisitedMin(v string) bool {
	if m.move1.valve == v || (m.move2 != nil && m.move2.valve == v) {
		return true
	}
	if m.prev == nil {
		return false
	}
	return m.prev.VisitedMin(v)
}

func (m *Minute) IsOpenMin(v string) bool {
	if (m.move1.valve == v && m.move1.openValve) || (m.move2 != nil && m.move2.valve == v && m.move2.openValve) {
		return true
	}
	if m.prev == nil {
		return false
	}
	return m.prev.IsOpenMin(v)
}

func (s *System) String() string {
	str := ""

	for name, connections := range s.connections {
		str += fmt.Sprintf("Valve %v (flow=%v) connects to %v\n", name, s.valves[name], strings.Join(connections, ", "))
	}
	return str
}

func (s *System) CandidateMoves(min *Minute, prev *Move) []*Move {
	candidates := []*Move{}

	if !min.IsOpenMin(prev.valve) {
		candidates = append(candidates, &Move{
			prev:          prev,
			minutesLeft:   prev.minutesLeft - 1,
			openValve:     true,
			valve:         prev.valve,
			valvePressure: s.valves[prev.valve],
		})
	}

	for _, conn := range s.connections[prev.valve] {
		// if prev.Visited(conn) {
		// 	continue
		// }

		candidates = append(candidates, &Move{
			prev:          prev,
			minutesLeft:   prev.minutesLeft - 1,
			openValve:     false,
			valve:         conn,
			valvePressure: s.valves[conn],
		})
	}

	return candidates
}

func parseInput(scanner *bufio.Scanner) (system System) {
	regex := regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)

	system.connections = make(map[string][]string)
	system.valves = make(map[string]int)

	for scanner.Scan() {
		matches := regex.FindStringSubmatch(scanner.Text())
		name := matches[1]
		rate := parseInt(matches[2])
		var connections []string

		for _, str := range strings.Split(matches[3], ",") {
			connections = append(connections, strings.TrimSpace(str))
		}

		system.valves[name] = rate

		system.connections[name] = connections
	}

	return system
}

func question1(scanner *bufio.Scanner) {
	system := parseInput(scanner)

	println(system.String())

	// path := system.GetPath("AA", "DD")

	// println(path.String())

	//greedy
	start := &Minute{
		move1: &Move{
			prev:          nil,
			minutesLeft:   30,
			valve:         "AA",
			openValve:     false,
			valvePressure: system.valves["AA"],
		},
		minutesLeft: 30,
	}

	finished := []*Minute{}
	current := []*Minute{start}
	var toContinue []*Minute
	moveCount := 30
	for len(current) > 0 {
		moveCount--
		toContinue = nil
		fmt.Printf("[%v] %v\n", moveCount, len(current))
		for _, m := range current {
			candidates := system.CandidateMoves(m, m.move1)
			if m.minutesLeft == 0 || len(candidates) == 0 {
				finished = append(finished, m)
				continue
			}
			for _, candidate := range candidates {
				nextMin := &Minute{
					prev:        m,
					minutesLeft: m.minutesLeft - 1,
					move1:       candidate,
				}
				toContinue = append(toContinue, nextMin)
			}
			// for _, candidate := range candidates {
			// 	// fmt.Printf("adding candidate from %v to %v (open=%v)\n", m.valve, candidate.valve, candidate.openValve)
			// 	toContinue = append(toContinue, candidate)
			// }
		}

		// for _, c := range toContinue {
		// 	println(c.String())
		// }

		if len(toContinue) > 0 && moveCount < 25 {

			sort.Slice(toContinue, func(i, j int) bool {
				vi := toContinue[i].cumulativePressureReleasedMin()
				vj := toContinue[j].cumulativePressureReleasedMin()
				return vi > vj
			})

			factor := math.Min((30.0-float64(moveCount))/30.0, 0.5)

			current = toContinue[:int(float64(len(toContinue))*factor)]
		} else {
			current = toContinue
		}

	}

	var best *Minute
	bestScore := -1
	for _, m := range finished {
		score := m.cumulativePressureReleasedMin()
		// println("-----------------------")
		// println(m.String())
		// println(score)
		// println("-----------------------")
		if score > bestScore {
			bestScore = score
			best = m
		}
	}

	// println(best.String())
	println(best.cumulativePressureReleasedMin())
}

func question2(scanner *bufio.Scanner) {
	system := parseInput(scanner)

	println(system.String())

	// path := system.GetPath("AA", "DD")

	// println(path.String())

	//greedy
	start := &Minute{
		move1: &Move{
			prev:          nil,
			minutesLeft:   26,
			valve:         "AA",
			openValve:     false,
			valvePressure: system.valves["AA"],
		},
		move2: &Move{
			prev:          nil,
			minutesLeft:   26,
			valve:         "AA",
			openValve:     false,
			valvePressure: system.valves["AA"],
		},
		minutesLeft: 26,
	}

	finished := []*Minute{}
	current := []*Minute{start}
	var toContinue []*Minute
	moveCount := 26
	for len(current) > 0 {
		moveCount--
		toContinue = nil
		fmt.Printf("[%v] %v\n", moveCount, len(current))
		for _, m := range current {
			candidates1 := system.CandidateMoves(m, m.move1)
			candidates2 := system.CandidateMoves(m, m.move2)
			if m.minutesLeft == 0 || len(candidates1) == 0 || len(candidates2) == 0 {
				finished = append(finished, m)
				continue
			}
			for _, candidate1 := range candidates1 {
				for _, candidate2 := range candidates2 {
					if candidate1.valve == candidate2.valve && candidate1.openValve && candidate2.openValve {
						// no shenanigans
						continue
					}
					nextMin := &Minute{
						prev:        m,
						minutesLeft: m.minutesLeft - 1,
						move1:       candidate1,
						move2:       candidate2,
					}
					toContinue = append(toContinue, nextMin)
				}
			}
			// for _, candidate := range candidates {
			// 	// fmt.Printf("adding candidate from %v to %v (open=%v)\n", m.valve, candidate.valve, candidate.openValve)
			// 	toContinue = append(toContinue, candidate)
			// }
		}

		// for _, c := range toContinue {
		// 	println(c.String())
		// }

		if len(toContinue) > 0 && moveCount < 22 {

			sort.Slice(toContinue, func(i, j int) bool {
				vi := toContinue[i].cumulativePressureReleasedMin()
				vj := toContinue[j].cumulativePressureReleasedMin()
				return vi > vj
			})

			factor := math.Min((30.0-float64(moveCount))/30.0, 0.7)

			current = toContinue[:int(float64(len(toContinue))*factor)]
		} else {
			current = toContinue
		}

	}

	var best *Minute
	bestScore := -1
	for _, m := range finished {
		score := m.cumulativePressureReleasedMin()
		// println("-----------------------")
		// println(m.String())
		// println(score)
		// println("-----------------------")
		if score > bestScore {
			bestScore = score
			best = m
		}
	}

	// println(best.String())
	println(best.cumulativePressureReleasedMin())

}

const isTest = true
const print = false

func main() {

	path := "./input.txt"
	if isTest {
		path = "./test.txt"
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if os.Args[1] == "1" {
		question1(scanner)
	} else {
		question2(scanner)
	}
}
