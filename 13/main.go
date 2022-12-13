package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

const print = false

type List struct {
	items []*List
	value *int
}

func (list *List) print() (str string) {
	if list.value == nil {
		str += "["
		for _, item := range list.items {
			str += item.print()
			str += ","
		}
		str += "]"
	} else {
		str += strconv.Itoa(*list.value)
	}
	return str
}

func (left *List) compare(right *List, indent string) int {
	if print {
		fmt.Printf(indent+"%v vs %v\n", left.print(), right.print())
	}
	if left.value != nil && right.value != nil {
		if *left.value-*right.value > 0 && print {
			fmt.Printf(indent+"\tLeft is bigger (bad) %v, %v\n", *left.value, *right.value)
		} else if *left.value-*right.value < 0 && print {
			fmt.Printf(indent+"\tLeft is smaller (good) %v, %v\n", *left.value, *right.value)
		}
		return *left.value - *right.value
	}
	if left.value == nil && right.value == nil {
		for i := 0; i < len(left.items); i++ {
			if i >= len(right.items) {
				if print {
					fmt.Println(indent + "Right is shorter (bad)")
				}
				return 1
			}

			v := left.items[i].compare(right.items[i], indent+"\t")
			if v != 0 {
				return v
			}
		}

		if len(left.items) < len(right.items) {
			if print {
				fmt.Println(indent + "Left is shorter (good)")
			}
			return -1
		}
		return 0
	}

	if left.value != nil {
		wrapped := List{
			items: []*List{
				{
					value: left.value,
				},
			},
		}
		return wrapped.compare(right, indent+"\t")
	}

	wrapped := List{
		items: []*List{
			{
				value: right.value,
			},
		},
	}
	return left.compare(&wrapped, indent+"\t")
}

type PacketPair struct {
	left  List
	right List
}

func (packet *PacketPair) print() {
	if print {
		fmt.Println(packet.left.print())
		fmt.Println(packet.right.print())
	}
}

func parseInt(str string) *int {
	v, err := strconv.Atoi(string(str))
	if err != nil {
		panic("Could not parse " + string(str))
	}
	return &v
}

func parseList(line string) List {
	var root List
	stack := []*List{}

	intBuff := ""

	flushIntBuffer := func() {
		if len(intBuff) > 0 {
			stack[len(stack)-1].items = append(stack[len(stack)-1].items, &List{items: []*List{}, value: parseInt(intBuff)})
			intBuff = ""
		}
	}

	for _, r := range line {
		switch r {
		case '[':
			{
				stack = append(stack, &List{
					items: []*List{},
					value: nil,
				})
			}
		case ']':
			{
				flushIntBuffer()
				if len(stack) > 1 {
					stack[len(stack)-2].items = append(stack[len(stack)-2].items, stack[len(stack)-1])
				}
				root = *stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		case ',':
			flushIntBuffer()
			continue
		default:
			{
				intBuff = intBuff + string(r)
			}
		}
	}
	return root
}

func parseInput(scanner *bufio.Scanner) (packets []PacketPair) {

	for {
		ok := scanner.Scan()
		if !ok {
			panic("Expected 1st packet")
		}
		line1 := scanner.Text()
		ok = scanner.Scan()
		if !ok {
			panic("Expected 2nd packet")
		}
		line2 := scanner.Text()

		packets = append(packets, PacketPair{
			left:  parseList(line1),
			right: parseList(line2),
		})

		ok = scanner.Scan()
		if !ok {
			break
		}
	}

	return packets
}

func unzip(pairs []PacketPair) (list []List) {
	for _, pair := range pairs {
		list = append(list, pair.left)
		list = append(list, pair.right)
	}
	return list
}

func question1(scanner *bufio.Scanner) {
	packets := parseInput(scanner)

	sum := 0
	for i, p := range packets {
		if p.left.compare(&p.right, "") < 0 {
			sum += i + 1
		}
	}

	fmt.Println(sum)
}

func question2(scanner *bufio.Scanner) {
	packets := unzip(parseInput(scanner))

	two := 2
	six := 6

	packet2 := List{
		items: []*List{
			{
				value: &six,
			},
		},
	}
	packet6 := List{
		items: []*List{
			{
				value: &two,
			},
		},
	}

	packets = append(packets, packet2)
	packets = append(packets, packet6)

	sort.Slice(packets, func(i, j int) bool {
		return packets[i].compare(&packets[j], "") < 0
	})

	ans := 1
	for i, p := range packets {
		if p.compare(&packet2, "") == 0 || p.compare(&packet6, "") == 0 {
			ans *= (i + 1)
		}
	}

	fmt.Println(ans)
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
