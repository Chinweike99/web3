package handlers

import (
	"cli-bank/internal/model"
	"fmt"
)


func ShowBalance(accounts []model.Account){
	fmt.Println("\nAccount Balances")
	for i, acc := range accounts{
		fmt.Printf("%d. %s: $%d\n", i+1, acc.Name, acc.Balance)
	}
}