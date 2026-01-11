package ui

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)


var reader = bufio.NewReader(os.Stdin)

func ReadLine() string{
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ReadInt(prompt string)(int, bool){
	fmt.Print(prompt)
	input := ReadLine()

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Enter a valid positive number")
		return 0, false
	}

	return value, true

}