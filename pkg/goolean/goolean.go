package goolean

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func BeginShell() (result string) {
	for {
		fmt.Print("> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()

		if err != nil {
			log.Fatal(err)
		}

		res := handleInput(scanner.Text())

		fmt.Printf("%s\n", res)
	}
}

func handleInput(input string) (result string) {
	return input
}
