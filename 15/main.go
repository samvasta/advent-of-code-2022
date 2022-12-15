package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const isTest = false
const print = false

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func parseInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		panic("wtf")
	}

	return v
}

type Sensor struct {
	x int
	y int

	beaconX int
	beaconY int

	dist int
}

func countNotSensors(y int, sensors []*Sensor) int {

	count := 0

	visited := make(map[int]bool)

	for _, sensor := range sensors {

		distToY := abs(y - sensor.y)

		if sensor.dist > distToY {
			start := sensor.x - abs(sensor.dist-distToY)
			end := sensor.x + abs(sensor.dist-distToY)
			for x := start; x <= end; x++ {
				if _, ok := visited[x]; !ok {
					visited[x] = true
					count++
				}
			}
		}
	}

	return count - 1 // obo b/c answer does not count beacons
}

func parseInput(scanner *bufio.Scanner) (sensors []*Sensor) {
	regex := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)

	for scanner.Scan() {
		matches := regex.FindStringSubmatch(scanner.Text())
		x := parseInt(matches[1])
		y := parseInt(matches[2])
		beaconX := parseInt(matches[3])
		beaconY := parseInt(matches[4])

		sensors = append(sensors, &Sensor{
			x:       x,
			y:       y,
			beaconX: beaconX,
			beaconY: beaconY,
			dist:    abs(x-beaconX) + abs(y-beaconY),
		})
	}

	return sensors
}

func question1(scanner *bufio.Scanner) {
	sensors := parseInput(scanner)

	row := 2_000_000
	if isTest {
		row = 10
	}
	fmt.Println(countNotSensors(row, sensors))
}

func question2(scanner *bufio.Scanner) {
	sensors := parseInput(scanner)

	limit := 4000000
	if isTest {
		limit = 20
	}
	toCheck := []uint64{}
	for _, sensor := range sensors {
		for y := sensor.y - sensor.dist; y < sensor.y+sensor.dist; y++ {
			if y < 0 || y > limit {
				continue
			}
			distToY := abs(y - sensor.y)
			diff := abs(distToY - sensor.dist)
			x1 := sensor.x + diff + 1
			x2 := sensor.x - diff - 1

			if x1 >= 0 && x1 <= limit {
				toCheck = append(toCheck, uint64(x1)<<32|uint64(y))
			}

			if x2 >= 0 && x2 <= limit {
				toCheck = append(toCheck, uint64(x2)<<32|uint64(y))
			}
		}

		if sensor.y-sensor.dist-1 >= 0 {
			toCheck = append(toCheck, uint64(sensor.x)<<32|uint64(sensor.y-sensor.dist-1))
		}

		if sensor.y+sensor.dist+1 <= limit {
			toCheck = append(toCheck, uint64(sensor.x)<<32|uint64(sensor.y+sensor.dist+1))
		}
	}

	for _, candidate := range toCheck {
		cx := int(candidate >> 32)
		cy := int(candidate & 0xffffff)
		isCovered := false
		for _, sensor := range sensors {
			if abs(sensor.x-cx)+abs(sensor.y-cy) <= sensor.dist {
				isCovered = true
				break
			}
		}

		if !isCovered {
			fmt.Printf("%v\n", cx*4000000+cy)
			break
		}
	}
}

func main() {

	path := "./input.txt"
	if isTest {
		path = "./test.txt"
	}

	file, err := os.Open(path)
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
