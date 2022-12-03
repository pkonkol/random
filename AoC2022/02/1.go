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
			"X": {"A": 1 + 3, "B": 1 + 0, "C": 1 + 6},
			"Y": {"A": 2 + 6, "B": 2 + 3, "C": 2 + 0},
			"Z": {"A": 3 + 0, "B": 3 + 6, "C": 3 + 3},
		}
		total += huj[us][enemy]
	}
	fmt.Printf("Total: %d\n", total)
}
