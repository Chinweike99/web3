package ui

import "fmt"

func ShowWelcome() {
	fmt.Println("=================================")
	fmt.Println(" Welcome to CLI Bank")
	fmt.Println("=================================")
}

func ShowMainMenu() {
	fmt.Println("\nMain Menu")
	fmt.Println("1. View Balances")
	fmt.Println("2. Transfer Money")
	fmt.Println("3. Exit")
	fmt.Print("Select an option: ")
}
