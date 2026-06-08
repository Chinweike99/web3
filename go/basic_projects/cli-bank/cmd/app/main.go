package main

import (
	"cli-bank/internal/handlers"
	"cli-bank/internal/model"
	"cli-bank/internal/ui"
	"fmt"
	"os"
)

func main() {
	accounts := []model.Account{
		{Name: "Checking", Balance: 1000},
		{Name: "Savings", Balance: 2000},
	}

	ui.ShowWelcome()
	for {
		ui.ShowMainMenu()
		choice := ui.ReadLine()

		switch choice {
		case "1":
			handlers.ShowBalance(accounts)
		case "2":
			handlers.Transfer(accounts)
		case "3":
			fmt.Println("\nThank you !...\n")
			os.Exit(0)
		default:
			fmt.Println("Invalid option.")
		}
	}

}