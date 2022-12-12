package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Cell struct {
	height rune

	x, y int

	left   *Cell
	right  *Cell
	top    *Cell
	bottom *Cell
}

func (cell *Cell) canStep(other *Cell) bool {
	if other == nil {
		return false
	}
	if other.height == 'a' && cell.height == 'S' {
		return true
	}

	if other.height == 'E' {
		return cell.height == 'z'
	}
	return other.height-cell.height <= 1
}

func (cell *Cell) hashPos() uint64 {
	return uint64(cell.x<<32) | uint64(cell.y)
}

func (cell *Cell) Equals(other *Cell) bool {
	return cell.x == other.x && cell.y == other.y
}

func (cell *Cell) DistSq(other *Cell) int {
	dx := cell.x - other.x
	dy := cell.y - other.y
	return dx*dx + dy*dy
}

func (cell *Cell) neighbors() (cells []*Cell) {
	if cell.canStep(cell.left) {
		cells = append(cells, cell.left)
	}

	if cell.canStep(cell.right) {
		cells = append(cells, cell.right)
	}
	if cell.canStep(cell.top) {
		cells = append(cells, cell.top)
	}
	if cell.canStep(cell.bottom) {
		cells = append(cells, cell.bottom)
	}
	return cells
}

type Grid struct {
	topLeft *Cell
	starts  []*Cell
	end     *Cell
	cells   []*Cell
}

type AStarState struct {
	prev    *AStarState
	current *Cell
	g       int
	h       float64
}

func (aStar *AStarState) F() float64 {
	return float64(aStar.g) + aStar.h
}

func (aStar *AStarState) Len() int {
	if aStar.prev == nil {
		return 0
	}
	return 1 + aStar.prev.Len()
}

func (aStar *AStarState) PrintPath() {
	fmt.Printf("%v (%v,%v)  ->  ", string(aStar.current.height), aStar.current.x, aStar.current.y)

	if aStar.prev == nil {
		return
	}
	aStar.prev.PrintPath()
}

func (grid *Grid) Print() {
	for firstOfRow := grid.topLeft; firstOfRow != nil; firstOfRow = firstOfRow.bottom {
		for current := firstOfRow; current != nil; current = current.right {
			fmt.Print(string(current.height))
		}
		fmt.Println()
	}
	fmt.Println()
}

func FindPath(from, to *Cell) *AStarState {
	visited := make(map[uint64]AStarState)
	toVisit := []AStarState{{
		prev:    nil,
		current: from,
		g:       0,
		h:       float64(from.DistSq(to)),
	}}

	visited[from.hashPos()] = toVisit[0]

	getNextState := func() AStarState {
		min := 99999999.0
		minIdx := 0
		for i, state := range toVisit {
			if state.F() < min {
				min = state.F()
				minIdx = i
			}
		}

		minState := toVisit[minIdx]
		toVisit = append(toVisit[:minIdx], toVisit[minIdx+1:]...)

		return minState
	}

	for len(toVisit) > 0 {
		next := getNextState()

		for _, neighbor := range next.current.neighbors() {
			cell := AStarState{
				current: neighbor,
				prev:    &next,
				g:       next.g + 1,
				h:       float64(neighbor.DistSq(to)) / float64(neighbor.height*neighbor.height),
			}

			if neighbor.Equals(to) {
				return &cell
			}

			if existing, ok := visited[neighbor.hashPos()]; !ok || existing.g > cell.g {
				toVisit = append(toVisit, cell)
				visited[neighbor.hashPos()] = cell
			}
		}
	}
	return nil
}

func parseInput(scanner *bufio.Scanner, startAtAllA bool) (grid Grid) {

	var lastRowStart *Cell = nil
	var lastRowCurrent *Cell = nil
	var lastCell *Cell = nil

	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		for x, height := range line {
			cell := new(Cell)
			cell.height = height
			cell.x = x
			cell.y = y
			cell.left = lastCell
			cell.top = lastRowCurrent

			if x == 0 && y == 0 {
				grid.topLeft = cell
			}

			if height == 'S' || (startAtAllA && height == 'a') {
				grid.starts = append(grid.starts, cell)
			} else if height == 'E' {
				grid.end = cell
			}

			if lastCell != nil {
				lastCell.right = cell
			}
			if lastRowCurrent != nil {
				lastRowCurrent.bottom = cell
				lastRowCurrent = lastRowCurrent.right
			}
			lastCell = cell
		}

		if lastRowStart == nil {
			lastRowStart = grid.topLeft
		} else {
			lastRowStart = lastRowStart.bottom
		}
		lastRowCurrent = lastRowStart
		lastCell = nil
		y++
	}

	return grid
}

func question1(scanner *bufio.Scanner) {
	grid := parseInput(scanner, false)

	path := FindPath(grid.starts[0], grid.end)

	fmt.Println(path.Len())
}

func question2(scanner *bufio.Scanner) {
	grid := parseInput(scanner, true)

	min := 999999999
	for _, start := range grid.starts {
		path := FindPath(start, grid.end)
		if path != nil {

			dist := path.Len()
			if dist < min {

				min = dist
			}
		}
	}

	fmt.Println(min)

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
