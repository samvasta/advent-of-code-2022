package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func isAllDifferent(chars []rune) bool {
	set := make(map[rune]bool)

	for _, c := range chars {
		if _, ok := set[c]; ok {
			return false
		}
		set[c] = true
	}

	return true
}

func findFirstMarker(input string, length int) (position int, before, rest string) {

	prev := make([]rune, length)

	prevIndex := 0

	for i, c := range input {
		prev[prevIndex] = c
		prevIndex++
		if prevIndex >= length {
			prevIndex = 0
		}

		if i >= 3 && isAllDifferent(prev) {
			position = i + 1
			before = input[:position]
			rest = input[position:]
			return
		}
	}

	return -1, input, ""
}

func question1(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		position, x, y := findFirstMarker(line, 4)

		fmt.Printf("%v, %v, %v\n", position, x, y)
	}
}

func question2(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		position, x, y := findFirstMarker(line, 14)

		fmt.Printf("%v, %v, %v\n", position, x, y)
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
