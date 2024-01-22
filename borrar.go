package main

import (
	"fmt"
	"strings"
)

func main() {

	command := "go run main.go"
	args := []string{"1", "2"}
	command_with_args := command + " " + strings.Join(args, " ")

	fmt.Println(command_with_args)
}
