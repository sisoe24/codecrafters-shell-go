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
	stdout  *os.File
	stderr  *os.File
	command string
	args    []string
}

func (c *Command) setOutput() {
	var args []string
	var redirect string
	var redirectArgs []string

	for i, arg := range c.args {
		if strings.Contains(arg, ">") {
			redirect = arg
			redirectArgs = c.args[i+1:]
			break
		} else {
			args = append(args, arg)
		}
	}

	c.args = args

	c.stdout = os.Stdout
	c.stderr = os.Stderr

	switch redirect {
	case ">", "1>":
		file, err := os.Create(redirectArgs[0])
		if err != nil {
			panic(err)
		}
		c.stdout = file
	}
}

func newCommand(input string) Command {
	args := parseArgs(strings.Split(input, " "))
	c := Command{command: args[0], args: args[1:]}
	c.setOutput()
	return c
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

		if char == ESCAPE {
			openEscape = true
			nextChar := text[i+1]

			if !openQuote {
				char = nextChar
				i++
			} else if !openSingleQuote {
				if nextChar == DOUBLE_QUOTE || nextChar == '$' || nextChar == ESCAPE {
					char = nextChar
					i++
				}
			}
		}

		if !openEscape && !openDoubleQuote && char == SINGLE_QUOTE {
			openSingleQuote = !openSingleQuote
			openQuote = !openQuote
			continue
		}

		if !openEscape && !openSingleQuote && char == DOUBLE_QUOTE {
			openDoubleQuote = !openDoubleQuote
			openQuote = !openQuote
			continue
		}

		if char == space && !openQuote {
			if (openEscape && item == "") || item != "" {
				items = append(items, item)
				item = ""
			}
		} else {
			item += string(char)
		}
		openEscape = false
	}

	if item != "" {
		items = append(items, item)
	}

	// fmt.Println("[debug] items:")
	// for _, item := range items {
	// 	fmt.Printf("  - %v\n", item)
	// }

	return items
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func binExists(paths string, bin string) string {
	for _, path := range strings.Split(paths, ":") {
		fp := filepath.Join(path, bin)
		if FileExists(fp) {
			return fmt.Sprintf("%s is %s/%s\n", bin, path, bin)
		}
	}
	return ""
}

func executeCommand(command Command, stdOut, stdErr *os.File) error {
	cmd := exec.Command(command.command, command.args...)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func parseFileArgs(c *Command) string {
	var args []string
	for _, arg := range c.args {
		if FileExists(arg) {
			args = append(args, arg)
		} else {
			fmt.Fprintf(c.stderr, "%s: %s: No such file or directory\n", c.command, arg)
		}
	}
	c.args = args
	return ""
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

		stdout := command.stdout
		stderr := command.stderr

		strArgs := strings.Join(command.args, " ")

		switch command.command {
		case "type":

			arg := command.args[0]

			if slices.Contains(builtins, arg) {
				fmt.Fprintf(stdout, "%s is a shell builtin\n", arg)
			} else if out := binExists(ENV_PATH, arg); out != "" {
				fmt.Fprint(stdout, out)
			} else {
				fmt.Fprintf(stderr, "%s: not found\n", arg)
			}

		case "cd":
			if strArgs == "~" {
				home, _ := os.UserHomeDir()
				os.Chdir(home)
			} else if FileExists(strArgs) {
				os.Chdir(strArgs)
			} else {
				fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", strArgs)
			}

		case "pwd":
			dir, _ := os.Getwd()
			fmt.Fprintf(stdout, "%s\n", dir)

		case "echo":
			fmt.Fprintf(stdout, "%s\n", strArgs)

		case "exit":
			os.Exit(0)

		case "ls", "cat":
			parseFileArgs(&command)
			fallthrough

		default:
			if err := executeCommand(command, stdout, stderr); err != nil {
				// technically not true but its good enough
				fmt.Fprintf(stdout, "%s: command not found\n", input)
			}
		}

		fmt.Fprint(os.Stdout, "$ ")

	}
}
