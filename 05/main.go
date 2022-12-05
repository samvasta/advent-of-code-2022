package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var regex, _ = regexp.Compile("move (\\d+) from (\\d+) to (\\d+)")

type Towers struct {
	stacks [][]rune
}

type Instruction struct {
	from  int
	to    int
	count int
}

func toInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}

func parseTowers(lines []string) (towers Towers) {
	for i := 0; i < len(lines[len(lines)-1]); i += 4 {
		towers.stacks = append(towers.stacks, []rune{})
	}

	for i := len(lines) - 2; i >= 0; i-- {
		for c := 1; c < len(lines[i]); c += 4 {
			if lines[i][c] == ' ' {
				continue
			}
			towers.stacks[(c-1)/4] = append(towers.stacks[(c-1)/4], (rune)(lines[i][c]))
		}
	}

	return towers
}

func parseInstructions(lines []string) (instructions []Instruction) {

	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		instructions = append(instructions, Instruction{from: toInt(matches[2]) - 1, to: toInt(matches[3]) - 1, count: toInt(matches[1])})
	}

	return instructions
}

func parseInput(scanner *bufio.Scanner) (towers Towers, instructions []Instruction) {
	var towersLines []string
	var instructionsLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		} else {

			towersLines = append(towersLines, line)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		} else {

			instructionsLines = append(instructionsLines, line)
		}
	}

	towers = parseTowers(towersLines)
	instructions = parseInstructions(instructionsLines)

	return towers, instructions
}

func applyOne(towers *Towers, instruction Instruction) {

	fromTower := towers.stacks[instruction.from]
	var move []rune

	if len(fromTower) > instruction.count {
		move = fromTower[len(fromTower)-instruction.count:]
		fromTower = fromTower[0 : len(fromTower)-instruction.count]
	} else {
		move = fromTower
		fromTower = []rune{}
	}

	towers.stacks[instruction.from] = fromTower

	for j := instruction.count - 1; j >= 0; j-- {
		towers.stacks[instruction.to] = append(towers.stacks[instruction.to], move[j])
	}

}

func applyOneV2(towers *Towers, instruction Instruction) {
	fromTower := towers.stacks[instruction.from]
	var move []rune

	if len(fromTower) > instruction.count {
		move = fromTower[len(fromTower)-instruction.count:]
		fromTower = fromTower[0 : len(fromTower)-instruction.count]
	} else {
		move = fromTower
		fromTower = []rune{}
	}

	towers.stacks[instruction.from] = fromTower

	for j := 0; j < instruction.count; j++ {
		towers.stacks[instruction.to] = append(towers.stacks[instruction.to], move[j])
	}

}

func applyAll(towers *Towers, instructions []Instruction, reverse bool) {
	for _, instruction := range instructions {
		if reverse {
			applyOne(towers, instruction)
		} else {
			applyOneV2(towers, instruction)
		}
	}
}

func question1(scanner *bufio.Scanner) {
	towers, instructions := parseInput(scanner)

	applyAll(&towers, instructions, true)

	for _, arr := range towers.stacks {
		fmt.Printf("%v", string(arr[len(arr)-1]))
	}
}

func question2(scanner *bufio.Scanner) {
	towers, instructions := parseInput(scanner)

	applyAll(&towers, instructions, false)

	for _, arr := range towers.stacks {
		fmt.Printf("%v", string(arr[len(arr)-1]))
	}

}

func main() {
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
