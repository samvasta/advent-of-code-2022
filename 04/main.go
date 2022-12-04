package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Section struct {
	Start int
	End   int
}

type SectionPair struct {
	A Section
	B Section
}

var regex, _ = regexp.Compile("(\\d+)-(\\d+),(\\d+)-(\\d+)")

func toInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}

func parsePair(line string) SectionPair {
	matches := regex.FindStringSubmatch(line)

	return SectionPair{
		A: Section{
			Start: toInt(matches[1]),
			End:   toInt(matches[2]),
		},
		B: Section{
			Start: toInt(matches[3]),
			End:   toInt(matches[4]),
		},
	}
}

func (this *SectionPair) Overlaps() bool {
	return (this.A.Start <= this.B.Start && this.A.End >= this.B.End) ||
		(this.A.Start >= this.B.Start && this.A.End <= this.B.End)
}

func (this *SectionPair) Overlaps2() bool {
	return (this.A.End >= this.B.Start) && (this.A.Start <= this.B.End)
}

func question1(scanner *bufio.Scanner) {
	count := 0
	for scanner.Scan() {
		pair := parsePair(scanner.Text())

		if pair.Overlaps() {
			count++
		}
	}
	fmt.Println(count)
}

func question2(scanner *bufio.Scanner) {
	count := 0
	for scanner.Scan() {
		pair := parsePair(scanner.Text())

		if pair.Overlaps2() {
			count++
		}
	}
	fmt.Println(count)

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
