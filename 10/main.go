package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type State uint8

const (
	noop State = iota
	addx1
	addx2
)

type CPU struct {
	reg   int
	state State
	cycle int

	history []int
}

var dot = color.New(color.FgGreen)
var sprite = color.New(color.FgBlack).Add(color.BgGreen)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (cpu *CPU) advanceProgramCounter() {
	if abs(((cpu.cycle-1)%40)-cpu.reg) <= 1 {
		sprite.Print("#")
	} else {
		dot.Print(".")
	}
	if cpu.cycle%40 == 0 {
		fmt.Println()
	}
	cpu.cycle++
	cpu.history = append(cpu.history, cpu.reg)
}

func execCmd(line string, cpu *CPU) {
	parts := strings.Split(line, " ")

	if parts[0] == "addx" {
		cpu.advanceProgramCounter()
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			panic("AAAHHH")
		}
		cpu.advanceProgramCounter()
		cpu.reg += num
	} else {
		cpu.advanceProgramCounter()
	}
}

func question1(scanner *bufio.Scanner) {
	cpu := CPU{
		reg:   1,
		state: noop,
		cycle: 1,
	}
	for scanner.Scan() {
		execCmd(scanner.Text(), &cpu)
	}

	sum := 0

	for i := 19; i < 230; i += 40 {
		sum += (cpu.history[i] * (i + 1))
		fmt.Printf("%v (%v * %v)\n", sum, cpu.history[i], i+1)
	}

	fmt.Println(sum)
}

func question2(scanner *bufio.Scanner) {
	// it was a 2-in-1 solution!
	question1(scanner)
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
