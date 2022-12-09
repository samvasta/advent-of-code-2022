package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pos struct {
	x int
	y int
}

func abs(v int) int {
	if v >= 0 {
		return v
	} else {
		return -v
	}
}
func normalize(v int) int {
	if v == 0 {
		return 0
	}
	return v / abs(v)
}

func (pos *Pos) moveToward(target *Pos, visitedMap *map[string]bool) {

	dx := target.x - pos.x
	dy := target.y - pos.y

	if abs(dx) > 1 || abs(dy) > 1 {
		pos.x += normalize(dx)
		pos.y += normalize(dy)
	}

	if visitedMap != nil {
		newPosHash := strconv.Itoa(pos.x) + strconv.Itoa(pos.y)
		(*visitedMap)[newPosHash] = true
	}

}

func (pos *Pos) applyInstruction(direction string) {
	switch direction {
	case "L":
		{
			pos.x--
			return
		}

	case "R":
		{
			pos.x++
			return
		}
	case "U":
		{
			pos.y--
			return
		}
	case "D":
		{
			pos.y++
			return
		}
	}
}

func question1(scanner *bufio.Scanner) {
	headPos := Pos{x: 0, y: 0}
	currentPos := Pos{x: 0, y: 0}
	visitedMap := make(map[string]bool)

	for scanner.Scan() {
		instruction := scanner.Text()
		parts := strings.Split(instruction, " ")

		count, err := strconv.Atoi(parts[1])

		if err != nil {
			panic("AAAAAHHHH")
		}
		for i := 0; i < count; i++ {
			headPos.applyInstruction(parts[0])
			currentPos.moveToward(&headPos, &visitedMap)
		}

	}

	fmt.Println(len(visitedMap))
}

func question2(scanner *bufio.Scanner) {

	nodes := []Pos{}

	for i := 0; i < 10; i++ {
		nodes = append(nodes, Pos{x: 0, y: 0})
	}
	visitedMap := make(map[string]bool)

	for scanner.Scan() {
		instruction := scanner.Text()
		parts := strings.Split(instruction, " ")

		count, err := strconv.Atoi(parts[1])

		if err != nil {
			panic("AAAAAHHHH")
		}
		for i := 0; i < count; i++ {
			nodes[0].applyInstruction(parts[0])
			for nodeIdx := 1; nodeIdx < len(nodes); nodeIdx++ {
				vm := &visitedMap
				if nodeIdx != len(nodes)-1 {
					vm = nil
				}
				nodes[nodeIdx].moveToward(&nodes[nodeIdx-1], vm)
			}
		}

	}

	fmt.Println(len(visitedMap))
}

func main() {

	// file, err := os.Open("./test.txt")
	// file, err := os.Open("./test2.txt")
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
