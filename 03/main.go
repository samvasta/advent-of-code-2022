package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func toPriority(thingy rune) int {
	if thingy < 91 {
		return int(thingy) - 38
	}
	return int(thingy) - 96
}

func convertToPriorities(line []rune) (a uint64) {
	a = 0

	for _, value := range line {
		a = a | (1 << toPriority(value))
	}

	return a
}

func splitIntoCompartments(line []rune) (a, b uint64) {
	aRaw := line[:len(line)/2]
	bRaw := line[len(line)/2:]

	return convertToPriorities(aRaw), convertToPriorities(bRaw)
}

func question1(scanner *bufio.Scanner) {
	sum := 0

	for scanner.Scan() {
		a, b := splitIntoCompartments([]rune(scanner.Text()))

		diff := a & b

		priority := 0

		for diff != 1 {
			diff = diff >> 1
			priority++
		}
		sum += priority
	}

	fmt.Println(sum)
}

func question2(scanner *bufio.Scanner) {
	sum := 0

	for scanner.Scan() {
		a := convertToPriorities([]rune(scanner.Text()))

		scanner.Scan()
		b := convertToPriorities([]rune(scanner.Text()))

		scanner.Scan()
		c := convertToPriorities([]rune(scanner.Text()))

		diff := a & b & c

		priority := 0

		for diff != 1 {
			diff = diff >> 1
			priority++
		}
		sum += priority
	}

	fmt.Println(sum)
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
