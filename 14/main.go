package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const print = false

func parseInt(str string) uint32 {
	v, err := strconv.Atoi(str)
	if err != nil {
		panic("wtf")
	}

	return uint32(v)
}

type Coord struct {
	x, y uint32
}

func (c *Coord) toHashedCoord() HashedCoord {
	return newHashedCoord(c.x, c.y)
}

type HashedCoord uint64

func newHashedCoord(x, y uint32) HashedCoord {
	return HashedCoord(uint64(x)<<32 | uint64(y))
}

type Tile struct {
	x uint32
	y uint32

	display rune
	isSolid bool
}

func (t *Tile) toHashedCoord() HashedCoord {
	return newHashedCoord(t.x, t.y)
}

type Simulation struct {
	sandSource HashedCoord

	sands []*Tile

	currentSand *Tile

	grid map[HashedCoord]*Tile

	minX uint32
	minY uint32
	maxX uint32
	maxY uint32

	isDone bool
}

func (sim *Simulation) ToString() string {
	s := ""

	for y := sim.minY; y <= sim.maxY+2; y++ {
		for x := sim.minX; x <= sim.maxX; x++ {
			tile, ok := sim.grid[newHashedCoord(x, y)]
			if ok && tile != nil {
				s += string(tile.display)
			} else {
				s += "."
			}
		}
		s += "\n"
	}

	return s
}

func (sim *Simulation) solid(coord HashedCoord) bool {
	t, ok := sim.grid[coord]
	if !ok || t == nil {
		return false
	}
	return t.isSolid
}

func (sim *Simulation) Step(floor bool) {
	if sim.currentSand == nil {
		panic("Current sand nil??")
	}

	bottom := newHashedCoord(sim.currentSand.x, sim.currentSand.y+1)

	if !sim.solid(bottom) {
		if print {
			fmt.Println("Moving down")
		}
		sim.grid[sim.currentSand.toHashedCoord()] = nil
		sim.currentSand.y++
		sim.grid[sim.currentSand.toHashedCoord()] = sim.currentSand

		if sim.currentSand.y > sim.maxY {
			if floor {
				sim.currentSand = nil
			} else {
				sim.isDone = true
			}
		}
		return
	}
	bottomLeft := newHashedCoord(sim.currentSand.x-1, sim.currentSand.y+1)
	if !sim.solid(bottomLeft) {
		if print {
			fmt.Println("Moving down-left")
		}
		sim.grid[sim.currentSand.toHashedCoord()] = nil
		sim.currentSand.y++
		sim.currentSand.x--
		sim.grid[sim.currentSand.toHashedCoord()] = sim.currentSand

		if sim.currentSand.y > sim.maxY && floor {
			sim.currentSand = nil
		}
		return
	}
	bottomRight := newHashedCoord(sim.currentSand.x+1, sim.currentSand.y+1)
	if !sim.solid(bottomRight) {
		if print {
			fmt.Println("Moving right")
		}
		sim.grid[sim.currentSand.toHashedCoord()] = nil
		sim.currentSand.y++
		sim.currentSand.x++
		sim.grid[sim.currentSand.toHashedCoord()] = sim.currentSand

		if sim.currentSand.y > sim.maxY && floor {
			sim.currentSand = nil
		}
		return
	}

	sim.currentSand = nil

}

func normalizedDiff(a, b uint32) int32 {
	if a > b {
		return -1
	} else if b > a {
		return 1
	}
	return 0
}

func parseCoords(line string) (coords []Coord) {
	parts := strings.Split(line, "->")

	for _, part := range parts {
		xy := strings.Split(strings.TrimSpace(part), ",")
		coords = append(coords, Coord{x: parseInt(xy[0]), y: parseInt(xy[1])})
	}

	return coords
}

func (sim *Simulation) addSolid(x, y uint32) {
	sim.grid[newHashedCoord(x, y)] = &Tile{
		x:       x,
		y:       y,
		display: '#',
		isSolid: true,
	}
	if x < sim.minX {
		sim.minX = x
	}
	if x > sim.maxX {
		sim.maxX = x
	}

	if y < sim.minY {
		sim.minY = y
	}
	if y > sim.maxY {
		sim.maxY = y
	}
}
func (sim *Simulation) line(from, to *Coord) {
	if print {
		fmt.Printf("Line from %v,%v to %v,%v\n", from.x, from.y, to.x, to.y)
	}
	if from.x != to.x {
		for x := int32(from.x); x != int32(to.x); x += normalizedDiff(from.x, to.x) {
			sim.addSolid(uint32(x), to.y)
		}
	} else {
		for y := int32(from.y); y != int32(to.y); y += normalizedDiff(from.y, to.y) {
			sim.addSolid(to.x, uint32(y))

		}
	}

	sim.addSolid(to.x, to.y)
}

func parseInput(scanner *bufio.Scanner) Simulation {

	simulation := Simulation{
		sandSource:  HashedCoord(0),
		sands:       []*Tile{},
		currentSand: nil,
		grid:        make(map[HashedCoord]*Tile),
		minX:        999999999,
		minY:        0,
		maxX:        0,
		maxY:        0,
	}

	curCoord := &Coord{}

	for scanner.Scan() {
		coords := parseCoords(scanner.Text())

		curCoord = &coords[0]
		for i := 1; i < len(coords); i++ {
			next := coords[i]
			simulation.line(curCoord, &next)
			curCoord = &next
		}
	}

	return simulation
}
func question1(scanner *bufio.Scanner) {
	sim := parseInput(scanner)

	if print {
		fmt.Println(sim.ToString())
	}

	for !sim.isDone {
		if sim.currentSand == nil {
			sand := Tile{
				x:       500,
				y:       0,
				display: 'o',
				isSolid: true,
			}
			sim.currentSand = &sand
			sim.grid[sand.toHashedCoord()] = &sand
			sim.sands = append(sim.sands, &sand)
		}

		sim.Step(false)
	}

	if print {
		fmt.Println("FINAL")
		fmt.Println(sim.ToString())
	}

	fmt.Println(len(sim.sands) - 1)
}

func question2(scanner *bufio.Scanner) {
	sim := parseInput(scanner)

	if print {
		fmt.Println(sim.ToString())
	}

	for !sim.isDone {
		if sim.currentSand == nil {
			sourceTile, ok := sim.grid[newHashedCoord(500, 0)]
			if !ok || sourceTile == nil {
				sand := Tile{
					x:       500,
					y:       0,
					display: 'o',
					isSolid: true,
				}
				sim.currentSand = &sand
				sim.grid[sand.toHashedCoord()] = &sand
				sim.sands = append(sim.sands, &sand)
			} else {
				break
			}
		}

		sim.Step(true)
	}

	if print {
		fmt.Println("FINAL")
		fmt.Println(sim.ToString())
	}

	fmt.Println(len(sim.sands))

}

func main() {
	// file, err := os.Open("./test.txt")
	file, err := os.Open("./input.txt")
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
