package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type Command struct {
	command string
	args    []string
}

func parseInput(input string) Command {
	args := strings.Split(input, " ")
	return Command{command: args[0], args: args[1:]}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func binExists(paths string, bin string) bool {
	for _, path := range strings.Split(paths, ":") {
		fp := filepath.Join(path, bin)
		if FileExists(fp) {
			fmt.Printf("%s is %s/%s\n", bin, path, bin)
			return true
		}
	}
	return false
}

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	scanner := bufio.NewScanner(os.Stdin)

	builtins := []string{
		"echo", "type", "exit",
	}

	ENV_PATH := os.Getenv("PATH")

	for scanner.Scan() {

		input := scanner.Text()
		command := parseInput(input)

		switch command.command {
		case "type":

			arg := command.args[0]

			if slices.Contains(builtins, arg) {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else if binExists(ENV_PATH, arg) {
        // kind of ugly having nothing to do here
			} else {
				fmt.Printf("%s: not found\n", arg)
			}

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
