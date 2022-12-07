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

var cmdRegex, _ = regexp.Compile("\\$ (ls|cd)\\s*(.*)")

type File struct {
	name string
	size int
}

type Dir struct {
	name   string
	size   int
	parent *Dir
	files  []*File
	dirs   []*Dir
}

func (dir *Dir) Print(depth int) {
	fmt.Printf("%v - %v (dir)\n", strings.Repeat("  ", depth), dir.name)

	for _, file := range dir.files {
		fmt.Printf("%v - %v (file, size=%v)\n", strings.Repeat("  ", depth+1), file.name, file.size)
	}

	for _, child := range dir.dirs {
		child.Print(depth + 1)
	}
}

func (dir *Dir) TotalSize() int {
	size := 0

	for _, file := range dir.files {
		size += file.size
	}

	for _, child := range dir.dirs {
		size += child.TotalSize()
	}

	dir.size = size

	return size
}

func (dir *Dir) CollectSmallDirs() []*Dir {
	smallDirs := []*Dir{}

	for _, child := range dir.dirs {
		childDirs := child.CollectSmallDirs()
		smallDirs = append(smallDirs, childDirs...)
		if child.size < 100000 {

			smallDirs = append(smallDirs, child)
		}

	}

	return smallDirs
}

func (dir *Dir) CollectDirs() []*Dir {
	dirs := []*Dir{}

	for _, child := range dir.dirs {
		childDirs := child.CollectDirs()
		dirs = append(dirs, childDirs...)
		dirs = append(dirs, child)

	}

	return dirs
}

func parseCDCmd(arg string, lines []string, currentDir *Dir) (newCurrentDir *Dir) {
	if arg == ".." {
		return currentDir.parent
	}

	for _, dir := range currentDir.dirs {
		if dir.name == arg {
			return dir
		}
	}

	panic("Could not find dir " + arg)

}

func parseLSCmd(lines []string, currentDir *Dir) (newCurrentDir *Dir) {
	for _, line := range lines {
		parts := strings.Split(line, " ")
		fmt.Printf("%v (%v, %v)\n", line, parts[0], parts[1])
		if parts[0] == "dir" {
			fmt.Println("adding dir " + parts[1])
			currentDir.dirs = append(currentDir.dirs, &Dir{
				name:   parts[1],
				size:   0,
				parent: currentDir,
				files:  []*File{},
				dirs:   []*Dir{},
			})
		} else {
			size, err := strconv.Atoi(parts[0])
			if err != nil {
				panic("AAAAAAHHHHH")
			}
			currentDir.files = append(currentDir.files, &File{
				name: parts[1],
				size: size,
			})
		}
	}

	return currentDir
}

func parseCmd(lines []string, currentDir *Dir) (newCurrentDir *Dir) {
	matches := cmdRegex.FindStringSubmatch(lines[0])
	if matches[1] == "ls" {
		return parseLSCmd(lines[1:], currentDir)
	} else if matches[1] == "cd" {
		return parseCDCmd(matches[2], lines[1:], currentDir)
	}
	panic("AAAAHHHH")
}

func parseAll(scanner *bufio.Scanner) (root Dir) {
	// skip 1st line b/c it's always "$ cd /"
	scanner.Scan()

	root = Dir{
		name:   "/",
		size:   0,
		parent: nil,
		files:  []*File{},
		dirs:   []*Dir{},
	}

	currentDir := &root

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '$' && len(lines) > 0 {
			currentDir = parseCmd(lines, currentDir)
			lines = lines[:0]
		}
		lines = append(lines, line)
	}
	currentDir = parseCmd(lines, currentDir)

	return root
}

func question1(scanner *bufio.Scanner) {
	root := parseAll(scanner)
	root.Print(0)
	root.TotalSize()
	smallDirs := root.CollectSmallDirs()

	sum := 0
	for _, dir := range smallDirs {
		sum += dir.size
	}

	fmt.Println(sum)
}

func question2(scanner *bufio.Scanner) {
	total := 70000000
	required := 30000000
	root := parseAll(scanner)
	root.Print(0)
	used := root.TotalSize()

	available := total - used
	toFree := required - available

	smallDirs := root.CollectDirs()

	smallest := total
	for _, dir := range smallDirs {
		if dir.size >= toFree && dir.size < smallest {
			smallest = dir.size
		}
	}

	fmt.Println(smallest)

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
