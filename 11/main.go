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

func parseInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		panic(str + " AAAAHHHH")
	}
	return v
}

func OpAdd(item int, operand int) int {
	return item + operand
}

func OpSubtract(item int, operand int) int {
	return item - operand
}

func OpMultiply(item int, operand int) int {
	return item * operand
}

func OpSquare(item int, operand int) int {
	return item * item
}

type Operation struct {
	operator    func(item int, operand int) int
	operatorStr string
	operand     int
}

func (operation *Operation) Apply(item int, modulo int) int {
	if modulo > 0 {
		return operation.operator(item, operation.operand) % modulo
	}
	return operation.operator(item, operation.operand)
}

type Monkey struct {
	items       []int
	operation   Operation
	testDivisor int
	trueMonkey  int
	falseMonkey int

	numInspections int
}

func (monkey *Monkey) Print() {
	fmt.Printf("Monkey (%v)\n\t%v\n\tnew = old %v %v\n\tdivisible by %v ? %v : %v\n", monkey.numInspections, monkey.items, monkey.operation.operatorStr, monkey.operation.operand, monkey.testDivisor, monkey.trueMonkey, monkey.falseMonkey)
}

func (monkey *Monkey) getNextMonkey(item int) int {
	if (item % monkey.testDivisor) == 0 {
		return monkey.trueMonkey
	}
	return monkey.falseMonkey
}

type Cohort struct {
	monkeys      []*Monkey
	commonModulo int
}

func (cohort *Cohort) SimulateOneRound(divByThree bool) (next Cohort) {
	next.commonModulo = cohort.commonModulo
	for _, monkey := range cohort.monkeys {
		next.monkeys = append(next.monkeys, &Monkey{
			items:          monkey.items[0:],
			operation:      monkey.operation,
			testDivisor:    monkey.testDivisor,
			trueMonkey:     monkey.trueMonkey,
			falseMonkey:    monkey.falseMonkey,
			numInspections: monkey.numInspections,
		})
	}

	for monkeyIdx, monkey := range next.monkeys {
		if len(monkey.items) == 0 {
			continue
		}

		next.monkeys[monkeyIdx].numInspections += len(monkey.items)

		for _, item := range monkey.items {
			intermediateValue := monkey.operation.Apply(item, cohort.commonModulo)
			newValue := intermediateValue
			if divByThree {
				newValue = newValue / 3
			}
			nextMonkeyIdx := monkey.getNextMonkey(newValue)
			next.monkeys[nextMonkeyIdx].items = append(next.monkeys[nextMonkeyIdx].items, newValue)

			// fmt.Printf("monkey %v applied %v %v %v and got %v (/3 = %v). Giving to %v\n", monkeyIdx, item, monkey.operation.operatorStr, monkey.operation.operand, intermediateValue, newValue, nextMonkeyIdx)
		}
		next.monkeys[monkeyIdx].items = next.monkeys[monkeyIdx].items[:0]

	}
	// for _, monkey := range next.monkeys {
	// 	monkey.Print()
	// }
	return next
}

func (cohort *Cohort) Simulate(divByThree bool, rounds int) Cohort {
	if rounds > 0 {
		next := cohort.SimulateOneRound(divByThree)
		return next.Simulate(divByThree, rounds-1)
	}
	return *cohort
}

func (monkey *Monkey) parseItems(line string) {
	r := regexp.MustCompile(`\d+`)
	for _, match := range r.FindAllString(line, -1) {
		monkey.items = append(monkey.items, parseInt(match))
	}
}
func getOperatorFn(opStr string) func(item int, operand int) int {
	switch opStr {
	case "-":
		return OpSubtract
	case "*":
		return OpMultiply
	case "+":
		return OpAdd
	}
	return OpAdd
}

func (monkey *Monkey) parseOperator(line string) {
	r := regexp.MustCompile(`\w*Operation: new = old (\+|-|\*) (\d+|old)`)
	matches := r.FindStringSubmatch(line)

	operator := getOperatorFn(matches[1])

	operand := 0
	if matches[2] == "old" {
		operator = OpSquare
	} else {
		operand = parseInt(matches[2])
	}

	monkey.operation = Operation{
		operator:    operator,
		operatorStr: matches[1],
		operand:     operand,
	}
}

func (monkey *Monkey) parseTest(lines []string) {
	testRegex := regexp.MustCompile(`\w*Test: divisible by (\d+)`)
	resultRegex := regexp.MustCompile(`\w*If (false|true): throw to monkey (\d+)`)

	monkey.testDivisor = parseInt(testRegex.FindStringSubmatch(lines[0])[1])

	monkey.trueMonkey = parseInt(resultRegex.FindStringSubmatch(lines[1])[2])
	monkey.falseMonkey = parseInt(resultRegex.FindStringSubmatch(lines[2])[2])

}

func parseMonkey(lines []string) *Monkey {
	monkey := Monkey{
		items: make([]int, 0),
		operation: Operation{
			operator:    OpAdd,
			operatorStr: "+",
			operand:     0,
		},
		testDivisor:    1,
		trueMonkey:     0,
		falseMonkey:    0,
		numInspections: 0,
	}
	monkey.parseItems(lines[1])
	monkey.parseOperator(lines[2])
	monkey.parseTest(lines[3:])

	return &monkey
}

func parseInput(scanner *bufio.Scanner) (cohort Cohort) {

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			cohort.monkeys = append(cohort.monkeys, parseMonkey(lines))
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}

	cohort.monkeys = append(cohort.monkeys, parseMonkey(lines))

	return cohort
}

func question1(scanner *bufio.Scanner) {
	cohort := parseInput(scanner)

	cohort.commonModulo = -1

	cohort = cohort.Simulate(true, 20)

	biggest := 0
	secondBiggest := 0

	for i, monkey := range cohort.monkeys {
		fmt.Printf("Monkey %v = %v items (%v)\n", i, monkey.numInspections, monkey.items)
		if monkey.numInspections > biggest {
			secondBiggest = biggest
			biggest = monkey.numInspections
		} else if monkey.numInspections > secondBiggest {
			secondBiggest = monkey.numInspections
		}
	}

	fmt.Printf("%v * %v = %v\n", biggest, secondBiggest, biggest*secondBiggest)
}

func question2(scanner *bufio.Scanner) {
	cohort := parseInput(scanner)

	cohort.commonModulo = 1
	for _, monkey := range cohort.monkeys {
		cohort.commonModulo *= monkey.testDivisor
	}

	cohort = cohort.Simulate(false, 10000)

	biggest := 0
	secondBiggest := 0

	for _, monkey := range cohort.monkeys {
		if monkey.numInspections > biggest {
			secondBiggest = biggest
			biggest = monkey.numInspections
		} else if monkey.numInspections > secondBiggest {
			secondBiggest = monkey.numInspections
		}
	}

	fmt.Printf("%v * %v = %v\n", biggest, secondBiggest, biggest*secondBiggest)
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
