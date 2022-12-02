package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func question1(scanner *bufio.Scanner) {
	maxSum := 0
	curSum := 0

	for scanner.Scan() {
		if scanner.Text() == "" {
			maxSum = max(maxSum, curSum)
			curSum = 0
		} else {
			value, err := strconv.Atoi(scanner.Text())

			if err != nil {
				log.Fatal(err)
			}

			curSum += value
		}
	}

	fmt.Println(maxSum)
}

func sumSlice(values []int) int {
	if len(values) == 1 {
		return values[0]
	}

	return values[0] + sumSlice(values[1:])
}

func question2(scanner *bufio.Scanner) {
	var sums []int
	curSum := 0

	for scanner.Scan() {
		if scanner.Text() == "" {
			sums = append(sums, curSum)
			curSum = 0
		} else {
			value, err := strconv.Atoi(scanner.Text())

			if err != nil {
				log.Fatal(err)
			}

			curSum += value
		}
	}

	sort.Slice(sums, func(i, j int) bool {
		return sums[i] > sums[j]
	})

	fmt.Println(sumSlice(sums[0:3]))
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
