package main

import (
	"bufio"
	"fmt"
	"os"
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

	var sum, max int
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			fmt.Printf("line empty: %s\n", t)
			sum = 0
		} else {
			kcal, err := strconv.Atoi(t)
			if err != nil {
				panic("wrong input")
			}
			sum += kcal
			fmt.Printf("line: %s, kcal: %d, sum: %d\n", t, kcal, sum)
		}
		if sum > max {
			max = sum
		}
		fmt.Printf("max is %d\n", max)
	}
}
