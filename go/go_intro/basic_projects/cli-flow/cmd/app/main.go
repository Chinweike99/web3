package main

import (
	"fmt"
	"os"

	"cli-flow/internal/ui"

)


func main() {
	ui.ShowWelcome()

	for {
		ui.ShowMainMenu()
		choice := ui.ReadLine()

		switch choice{
		case "1":
			handleBalance()
		case "2":
			handleTransfer()
		case "3":
			fmt.Println("\nGoodbye")
			os.Exit(0)
		default:
			fmt.Println("\nInvalid option, please try again.")
		}
	}
}


func handleBalance(){
	fmt.Println("\n Your balance is $1000")
}

func handleTransfer(){
	fmt.Println("\nEnter amount to transfer")
	amount := ui.ReadLine()

	fmt.Printf("Transfer $%s ....\n", amount)
	fmt.Println("Transfer was successful")
}