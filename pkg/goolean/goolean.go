package goolean

import (
	"fmt"
)

func BeginShell() (result string) {
	for {
		input := ""

		fmt.Print("> ")

		fmt.Scanf("%s", &input)

		res := handleInput(&input)

		fmt.Println("%s", res)
	}
}

func handleInput(input *string) (result string) {
	return ""
}
