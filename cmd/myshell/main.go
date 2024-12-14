package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type Command struct {
  command string
  args []string
}


func parseInput(input string) Command {
  args := strings.Split(input, " ")
  return Command{command: args[0], args: args[1:]}
}

func main() {
  fmt.Fprint(os.Stdout, "$ ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

    input := scanner.Text()
    command := parseInput(input)

		switch  command.command {
    case "echo":
			fmt.Printf("%s\n", strings.Join(command.args, " "))
    case "exit":
      os.Exit(0)
		default:
			fmt.Printf("%s: command not found\n", input)
		}

    fmt.Fprint(os.Stdout, "$ ")
	}
}
