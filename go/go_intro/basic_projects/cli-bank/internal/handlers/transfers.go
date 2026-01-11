package handlers

import (
	"cli-bank/internal/model"
	"cli-bank/internal/ui"
	"fmt"
)


func Transfer(accounts []model.Account){
	ShowAccounts(accounts)

	fmt.Print("\nFrom account (or B to go back): ")
	fromInput := ui.ReadLine()
	if fromInput == "B" || fromInput == "b" {
		return
	}

	fmt.Print("To account (or B to go back): ")
	toInput := ui.ReadLine()
	if toInput == "B" || toInput == "b" {
		return
	}

	from, to := parseIndex(fromInput), parseIndex(toInput)
	if from < 0 || to < 0 || from == to || 
		from >= len(accounts) || to >= len(accounts) {
		fmt.Println("Invalid account selection")
		return
	}

	amount, ok := ui.ReadInt("Enter amount t transfer: ")
	if !ok || accounts[from].Balance < amount {
		fmt.Println("Innsufficient balance 🥲😏")
		return
	}

	fmt.Printf("\nConfirm transfer of $%d from %s to %s? (y/n): ", amount, accounts[from].Name, accounts[to].Name)

	confirm := ui.ReadLine()
	if confirm != "y"{
		fmt.Println("Transfer cancelled")
		return
	}

	accounts[from].Balance -= amount
	accounts[to].Balance += amount
	fmt.Println("Transfer Successful")
}


func ShowAccounts(accounts []model.Account){
	fmt.Println("\nAccounts:")
	for i, acc := range accounts {
		fmt.Printf("%d. %s ($%d)\n", i+1, acc.Name, acc.Balance)
	}
}


func parseIndex(input string) int {
	var i int
	_, err := fmt.Sscanf(input, "%d", &i)
	if err != nil {
		return -1
	}
	return i - 1
}
