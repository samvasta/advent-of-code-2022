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
			fmt.Printf("found duplicate %v\n", c)
			return false
		}
		set[c] = true
	}

	return true
}

func findFirstPacketMarker(input string) (position int, before, rest string) {

	prev := []rune{0, 0, 0, 0}

	prevIndex := 0

	for i, c := range input {
		prev[prevIndex] = c
		prevIndex++
		if prevIndex >= 4 {
			prevIndex = 0
		}

		fmt.Println(prev)
		if i >= 3 && isAllDifferent(prev) {
			position = i + 1
			before = input[:position]
			rest = input[position:]
			return
		}
	}

	return -1, input, ""
}

func findFirstMessageMarker(input string) (position int, before, rest string) {

	prev := []rune{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	prevIndex := 0

	for i, c := range input {
		prev[prevIndex] = c
		prevIndex++
		if prevIndex >= 14 {
			prevIndex = 0
		}

		fmt.Println(prev)
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
		position, x, y := findFirstPacketMarker(line)

		fmt.Printf("%v, %v, %v\n", position, x, y)
	}
}

func question2(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		position, x, y := findFirstMessageMarker(line)

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
