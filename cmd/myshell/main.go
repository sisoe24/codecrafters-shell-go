package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

const (
	space        = ' '
	SINGLE_QUOTE = '\''
	DOUBLE_QUOTE = '"'
	ESCAPE       = '\\'
)

func trimQuotes(text string) string {
	for _, quote := range []string{"'", "\""} {
		text = strings.Trim(text, quote)
	}
	return text
}

func parseArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	var item string
	var items []string

	text := strings.Join(args, " ")

	openQuote := false
	openSingleQuote := false
	openDoubleQuote := false
	openEscape := false

	for i := 0; i < len(text); i++ {
		char := text[i]

		if char == ESCAPE && !openSingleQuote {
			openEscape = true
			i++
			char = text[i]
		}

		if !openEscape && !openDoubleQuote && char == SINGLE_QUOTE {
			openSingleQuote = !openSingleQuote
			openQuote = !openQuote
		}

		if !openEscape && char == DOUBLE_QUOTE {
			openDoubleQuote = !openDoubleQuote
			openQuote = !openQuote
		}

		if char == space && !openQuote {

			if (openEscape && item == "") || item != "" {
				items = append(items, trimQuotes(item))
				item = ""
			}

		} else {
			if openEscape && openDoubleQuote  {
				item += string('\\') + string(char)
			} else {
				item += string(char)
			}
		}
		openEscape = false
	}

	if item != "" {
		items = append(items, trimQuotes(item))
	}

	// fmt.Println("[debug] items:")
	// for _, item := range items {
	// 	fmt.Printf("  - %v\n", item)
	// }

	return items
}

func newCommand(input string) Command {
	args := strings.Split(input, " ")
	return Command{command: args[0], args: parseArgs(args[1:])}
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

func executeCommand(command Command) error {
	cmd := exec.Command(command.command, command.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	scanner := bufio.NewScanner(os.Stdin)

	builtins := []string{
		"echo", "type", "exit", "pwd", "cd",
	}

	ENV_PATH := os.Getenv("PATH")

	for scanner.Scan() {

		input := scanner.Text()
		command := newCommand(input)
		strArgs := strings.Join(command.args, " ")

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

		case "cd":
			if strArgs == "~" {
				home, _ := os.UserHomeDir()
				os.Chdir(home)
			} else if FileExists(strArgs) {
				os.Chdir(strArgs)
			} else {
				fmt.Printf("cd: %s: No such file or directory\n", strArgs)
			}

		case "pwd":
			dir, _ := os.Getwd()
			fmt.Printf("%s\n", dir)
		case "echo":
			fmt.Printf("%s\n", strArgs)
		case "exit":
			os.Exit(0)
		default:
			if err := executeCommand(command); err != nil {
				// technically not true but its good enough
				fmt.Printf("%s: command not found\n", input)
			}
		}

		fmt.Fprint(os.Stdout, "$ ")
	}
}
