package main

import (
	"bufio"
	"fmt"
	"os"
  "strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input

	reader := bufio.NewReader(os.Stdin)
	command, _ := reader.ReadString('\n')

  command = strings.Trim(command, "\n")

	switch command {
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}
