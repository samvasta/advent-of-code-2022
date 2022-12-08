package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type Tree struct {
	height int
}

type Grid struct {
	trees    [][]Tree
	rowCount int
	colCount int
}

func (grid *Grid) calcCell(row, col int) (isVisible bool, viewingDist int) {
	if row == 0 ||
		col == 0 ||
		row == grid.rowCount-1 ||
		col == grid.colCount-1 {
		return true, 0
	}

	viewingDist = 1
	isVisible = false

	targetHeight := grid.trees[row][col].height

	for r := row - 1; r >= 0; r-- {
		height := grid.trees[r][col].height
		if r == 0 && height < targetHeight {
			isVisible = true
			viewingDist *= row
		}

		if height >= targetHeight {
			viewingDist *= (row - r)
			break
		}
	}

	for r := row + 1; r <= grid.rowCount; r++ {
		if r == grid.rowCount {
			isVisible = true
			viewingDist *= (r - row - 1)
			break
		}

		height := grid.trees[r][col].height
		if height >= targetHeight {
			viewingDist *= (r - row)
			break
		}
	}

	for c := col - 1; c >= 0; c-- {
		height := grid.trees[row][c].height
		if c == 0 && height < targetHeight {
			isVisible = true
			viewingDist *= col
		}

		if height >= targetHeight {
			viewingDist *= (col - c)
			break
		}
	}

	for c := col + 1; c <= grid.colCount; c++ {
		if c == grid.colCount {
			isVisible = true
			viewingDist *= (c - col - 1)
			break
		}

		height := grid.trees[row][c].height
		if height >= targetHeight {
			viewingDist *= (c - col)
			break
		}
	}

	return isVisible, viewingDist
}

func parseGrid(scanner *bufio.Scanner) (grid Grid) {
	grid = Grid{
		trees:    [][]Tree{},
		rowCount: 0,
		colCount: 0,
	}

	for scanner.Scan() {
		treesStr := scanner.Text()

		trees := []Tree{}
		for _, t := range treesStr {
			height, _ := strconv.Atoi(string(t))
			trees = append(trees, Tree{
				height: height,
			})
		}

		grid.trees = append(grid.trees, trees)
	}

	grid.rowCount = len(grid.trees)
	grid.colCount = len(grid.trees[0])

	return grid
}

func question1(scanner *bufio.Scanner) {
	grid := parseGrid(scanner)

	light := color.New(color.FgHiGreen)
	dark := color.New(color.FgHiBlack)
	numVisible := 0
	for r := 0; r < grid.rowCount; r++ {
		for c := 0; c < grid.colCount; c++ {
			isVisible, _ := grid.calcCell(r, c)
			if isVisible {
				numVisible += 1
				light.Print(grid.trees[r][c].height)
			} else {
				dark.Print(grid.trees[r][c].height)
			}
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(numVisible)

}

func question2(scanner *bufio.Scanner) {
	grid := parseGrid(scanner)

	light := color.New(color.FgHiGreen)
	dark := color.New(color.FgHiBlack)
	bestScore := 0
	for r := 0; r < grid.rowCount; r++ {
		for c := 0; c < grid.colCount; c++ {
			isVisible, dist := grid.calcCell(r, c)

			if isVisible {
				if dist > bestScore {
					bestScore = dist
				}
				light.Print(dist)
			} else {
				dark.Print(dist)
			}
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(bestScore)
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
