package main

import "fmt"

func makeCounter() func() int {
	count := 0

	return func() int {
		count++
		return count
	}
}


func main_() {
	counterA := makeCounter()
	counterB := makeCounter()

	fmt.Println(counterA()) 
	fmt.Println(counterA())

	fmt.Println(counterB())
}
