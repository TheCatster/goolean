package goolean

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	/*
	 * Sample input: (NOT A AND NOT B) OR (NOT A OR NOT B)
	 * Expected output: NOT A OR NOT B
	 * Truth Table:
	 */
	data := strings.Split(input, " ")
	tokens := tokenize(data)

	return strings.Join(tokens, ", ")
}

func tokenize(data []string) []string {
	var tokens []string

	for i := range data {
		p := data[i]
		if strings.Contains(p, "(") {
			cnt := strings.Join(strings.Split(p, "("), "")
			tokens = append(tokens, "(")
			tokens = append(tokens, cnt)
		} else if strings.Contains(p, ")") {
			cnt := strings.Join(strings.Split(p, ")"), "")
			tokens = append(tokens, cnt)
			tokens = append(tokens, ")")
		} else {
			tokens = append(tokens, p)
		}
	}

	return tokens
}
