package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func buildLookupQ1() map[string]map[string]int {
	var scores = make(map[string]map[string]int)
	scores["A"] = make(map[string]int)
	scores["B"] = make(map[string]int)
	scores["C"] = make(map[string]int)

	scores["A"]["X"] = 4 // R R
	scores["A"]["Y"] = 8 // R P
	scores["A"]["Z"] = 3 // R S
	scores["B"]["X"] = 1 // P R
	scores["B"]["Y"] = 5 // P P
	scores["B"]["Z"] = 9 // P S
	scores["C"]["X"] = 7 // S R
	scores["C"]["Y"] = 2 // S P
	scores["C"]["Z"] = 6 // S S

	return scores
}

func buildLookupQ2() map[string]map[string]int {
	var scores = make(map[string]map[string]int)
	scores["A"] = make(map[string]int)
	scores["B"] = make(map[string]int)
	scores["C"] = make(map[string]int)

	scores["A"]["X"] = 3 // R lose
	scores["A"]["Y"] = 4 // R draw
	scores["A"]["Z"] = 8 // R win
	scores["B"]["X"] = 1 // P lose
	scores["B"]["Y"] = 5 // P draw
	scores["B"]["Z"] = 9 // P win
	scores["C"]["X"] = 2 // S lose
	scores["C"]["Y"] = 6 // S draw
	scores["C"]["Z"] = 7 // S win

	return scores
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lookup map[string]map[string]int
	if os.Args[1] == "1" {
		lookup = buildLookupQ1()
	} else {
		lookup = buildLookupQ2()
	}

	sum := 0

	for scanner.Scan() {
		moves := strings.Split(scanner.Text(), " ")

		score := lookup[moves[0]][moves[1]]

		sum += score
	}

	fmt.Println(sum)
}
