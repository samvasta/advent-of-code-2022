package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

type Rock struct {
	height   int
	width    int
	pattern  []uint8
	nextRock func() *Rock
}

func HLine() *Rock {
	return &Rock{
		height: 1,
		width:  4,
		pattern: []uint8{
			0b11110000,
		},
		nextRock: Plus,
	}
}

func Plus() *Rock {
	return &Rock{
		height: 3,
		width:  3,
		pattern: []uint8{
			0b01000000,
			0b11100000,
			0b01000000,
		},
		nextRock: ReverseL,
	}
}

func ReverseL() *Rock {
	return &Rock{
		height: 3,
		width:  3,
		pattern: []uint8{
			0b11100000,
			0b00100000,
			0b00100000,
		},
		nextRock: VLine,
	}
}
func VLine() *Rock {
	return &Rock{
		height: 4,
		width:  1,
		pattern: []uint8{
			0b10000000,
			0b10000000,
			0b10000000,
			0b10000000,
		},
		nextRock: Square,
	}
}

func Square() *Rock {
	return &Rock{
		height: 2,
		width:  2,
		pattern: []uint8{
			0b11000000,
			0b11000000,
		},
		nextRock: HLine,
	}
}

type State = []uint8

type Sim struct {
	state *State

	numSettledRocks            uint64
	currentRock                *Rock
	currentRockX, currentRockY uint64

	maxHeight uint64
}

func StateString(state *State) string {
	lines := []string{}

	for _, line := range *state {
		lines = append(lines, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%08b", line), "0", "."), "1", "#"))
	}

	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return strings.Join(lines, "\n") + "\n-------\n"
}

func (s *Sim) String() string {
	return StateString(s.state)
}

func EraseRock(rock *Rock, x, y uint64, state *State) *State {
	next := append(State(nil), *state...)

	for i := 0; i < rock.height; i++ {
		if i+int(y) >= len(next) {
			continue
		}
		next[i+int(y)] = next[i+int(y)] & (^(rock.pattern[i] >> x))
	}

	return &next
}

func DoesRockFit(rock *Rock, x, y uint64, state *State) bool {

	if y < 0 {
		if print {
			fmt.Println("no fit - y < 0")
		}
		return false
	}

	for i := 0; i < rock.height; i++ {
		if i+int(y) >= len(*state) {
			continue
		}

		if int(x)+rock.width > 7 {
			if print {
				fmt.Printf("no fit - x > 7 (%v)\n", int(x)+rock.width)
			}
			return false
		}

		result := (*state)[i+int(y)] & ((rock.pattern[i]) >> int(x))
		if result > 0 {
			if print {
				fmt.Printf("no fit - result != 0 (%08b -> %08b)\n", (*state)[i+int(y)], result)
			}
			return false
		}
	}

	return true
}

func AddRock(rock *Rock, x, y uint64, state *State) *State {
	next := append(State(nil), *state...)

	for i := 0; i < rock.height; i++ {
		for i+int(y) >= len(next) {
			next = append(next, 0)
		}

		if print {
			fmt.Printf("%v: %08b => %08b\n", i+int(y), next[i+int(y)], next[i+int(y)]|((rock.pattern[i])>>x))
		}

		next[i+int(y)] = next[i+int(y)] | ((rock.pattern[i]) >> x)
	}

	return &next
}

func (s *Sim) MoveRockRight() {
	if print {
		fmt.Println("Moving right")
	}
	if s.currentRockX+uint64(s.currentRock.width) > 7 {
		return
	}
	without := EraseRock(s.currentRock, s.currentRockX, s.currentRockY, s.state)
	if DoesRockFit(s.currentRock, s.currentRockX+1, s.currentRockY, without) {
		s.state = AddRock(s.currentRock, s.currentRockX+1, s.currentRockY, without)
		s.currentRockX++
	}
}

func (s *Sim) MoveRockLeft() {
	if print {
		fmt.Println("Moving left")
	}
	if s.currentRockX == 0 {
		return
	}
	without := EraseRock(s.currentRock, s.currentRockX, s.currentRockY, s.state)
	if DoesRockFit(s.currentRock, s.currentRockX-1, s.currentRockY, without) {
		s.state = AddRock(s.currentRock, s.currentRockX-1, s.currentRockY, without)
		s.currentRockX--
	}

}
func (s *Sim) MoveRockDown() (success bool) {
	if print {
		fmt.Printf("Moving down %v\n", s.currentRockY)
	}

	without := EraseRock(s.currentRock, s.currentRockX, s.currentRockY, s.state)
	if s.currentRockY > 0 && DoesRockFit(s.currentRock, s.currentRockX, s.currentRockY-1, without) {
		s.state = AddRock(s.currentRock, s.currentRockX, s.currentRockY-1, without)
		s.currentRockY--
		return true
	}

	currentRockHeight := uint64(s.currentRockY + uint64(s.currentRock.height) - 1)

	if currentRockHeight > s.maxHeight {
		s.maxHeight = currentRockHeight
	}
	return false
}

func (s *Sim) Move(left bool) {
	if left {
		s.MoveRockLeft()
	} else {
		s.MoveRockRight()
	}

	ok := s.MoveRockDown()
	if !ok {
		if print {
			fmt.Println("finished rock. moving to next")
		}

		s.currentRock = s.currentRock.nextRock()
		s.currentRockX = 2
		s.currentRockY = s.maxHeight + 4
		s.state = AddRock(s.currentRock, s.currentRockX, s.currentRockY, s.state)
		s.numSettledRocks++
	}
}

func parseInput(scanner *bufio.Scanner) []bool {
	scanner.Scan()

	moves := []bool{}

	for _, c := range scanner.Text() {
		if c == '<' {
			moves = append(moves, true)
		} else if c == '>' {
			moves = append(moves, false)
		} else {
			panic("AAAHHH " + string(c))
		}
	}

	return moves
}

func question1(scanner *bufio.Scanner) {
	sim := Sim{
		state:        &[]uint8{},
		currentRock:  HLine(),
		currentRockX: 2,
		currentRockY: 3,
		maxHeight:    0,
	}

	moves := parseInput(scanner)

	for sim.numSettledRocks < 2022 {
		for _, moveLeft := range moves {
			sim.Move(moveLeft)
			if print {
				println(sim.String())
				println("\n")
			}
			if sim.numSettledRocks == 2022 {
				break
			}
		}
	}

	fmt.Println(sim.maxHeight)
}

func question2(scanner *bufio.Scanner) {
	sim := Sim{
		state:        &[]uint8{},
		currentRock:  HLine(),
		currentRockX: 2,
		currentRockY: 3,
		maxHeight:    0,
	}

	moves := parseInput(scanner)

	for sim.numSettledRocks < 1000000000000 {
		for _, moveLeft := range moves {
			sim.Move(moveLeft)
			if print {
				println(sim.String())
				println("\n")
			}
			if sim.numSettledRocks == 1000000000000 {
				break
			}
			if sim.numSettledRocks%100000000000 == 0 {
				fmt.Printf("placed %v so far\n", sim.numSettledRocks)
			}
		}
	}

	fmt.Println(sim.maxHeight)
}

var print = false

func main() {

	testPtr := flag.Bool("t", false, "use test input")
	flag.BoolVar(&print, "p", false, "print stuff")

	flag.Parse()

	fmt.Printf("test=%v, print=%v\n", *testPtr, print)

	path := "./input.txt"
	if *testPtr {
		path = "./test.txt"
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if flag.Args()[0] == "1" {
		question1(scanner)
	} else {
		question2(scanner)
	}
}
