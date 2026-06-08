package ui

import "fmt"

func ShowWelcome(){
	fmt.Println("============================================")
	fmt.Println(" Welcome to CLI learning flow with go")
	fmt.Println("=============================================")
}


func ShowMainMenu(){
	fmt.Println("\nWhat do you want to do ?")
	fmt.Println("1. Check Balance")
	fmt.Println("2. Make Transfer")
	fmt.Println("3. Exit")
	fmt.Print("Select an option: ")
}