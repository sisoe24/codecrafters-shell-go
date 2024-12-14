package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
  fmt.Fprint(os.Stdout, "$ ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

    command := scanner.Text()

		switch command {
    case "exit 0":
      os.Exit(0)
		default:
			fmt.Printf("%s: command not found\n", command)
		}

    fmt.Fprint(os.Stdout, "$ ")
	}
}
