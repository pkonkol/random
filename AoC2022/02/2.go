package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("asdf\n")
	file, err := os.Open("input")
	if err != nil {
		panic("wrong filename")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var total int
	for scanner.Scan() {
		enemy := scanner.Text()
		scanner.Scan()
		us := scanner.Text()
		fmt.Printf("enemy: %s us: %s\n", enemy, us)
		huj := map[string]map[string]int{
			"X": {"A": 3, "B": 1, "C": 2},
			"Y": {"A": 3 + 1, "B": 3 + 2, "C": 3 + 3},
			"Z": {"A": 6 + 2, "B": 6 + 3, "C": 6 + 1},
		}
		total += huj[us][enemy]
	}
	fmt.Printf("Total: %d\n", total)
}
