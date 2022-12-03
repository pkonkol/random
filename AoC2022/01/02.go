package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	fmt.Printf("asdf\n")
	file, err := os.Open("input01")
	if err != nil {
		panic("wrong filename")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var sum int
	var kcals []int
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			fmt.Printf("line empty: %s\n", t)
			kcals = append(kcals, sum)
			sum = 0
		} else {
			kcal, err := strconv.Atoi(t)
			if err != nil {
				panic("wrong input")
			}
			sum += kcal
			fmt.Printf("line: %s, kcal: %d, sum: %d\n", t, kcal, sum)
		}
	}
	sort.Ints(kcals)
	l := len(kcals) - 1
	fmt.Printf("-1: %d, -2: %d, -3: %d, total: %d\n", kcals[l], kcals[l-1], kcals[l-2], kcals[l]+kcals[l-1]+kcals[l-2])
}
