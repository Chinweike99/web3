package ui

import (
	"bufio"
	"os"
	"strings"
)


var reader = bufio.NewReader(os.Stdin)

func ReadLine() string{
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}