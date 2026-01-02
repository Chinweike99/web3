package main

import "fmt"

// func main() {
// 	s := "Hello Chinweike"
// 	fmt.Println(s)
// }


func main_closures(){
	multiplyByTwo := multiplyBy(2)
	multiplyByFour := multiplyBy(4)

	result1 := multiplyByTwo(5)
	result2 := multiplyByFour(6)

	fmt.Println(result1, result2)
	arithmetic()
}


func arithmetic() {
	a := 10
	b := 3

	fmt.Println("Addition:", a+b)
	fmt.Println("Subtraction:", a-b)
	fmt.Println("Multiplication:", a*b)
	fmt.Println("Division:", a/b) 
	fmt.Println("Remainder:", a%b)
}


func multiplyBy(multiplier int) func(int) int{
	return func(i int) int{
		return i * multiplier
	}
}


