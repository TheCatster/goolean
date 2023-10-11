package goolean

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"
)

type Token struct {
	Type  int
	Value string
}

type Node struct {
	Left, Right *Node
	Value       Token
}

const (
	OPERATOR = iota
	VARIABLE
)

func Run() {
	runRepl()
}

func runRepl() {
	for {
		input := getLine()
		if input == "exit" {
			return
		} else {
			tree, err := parse(input)

			if err != nil {
				fmt.Println(err)
				continue
			}

			simplifiedTree := simplify(tree)
			simplifiedExprStr := printExpr(simplifiedTree)
			fmt.Println(simplifiedExprStr)
		}
	}
}

func getLine() string {
	fmt.Print("goolean> ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()

	if err != nil {
		log.Fatal(err)
	}

	return scanner.Text()
}

func isOperator(ch string) bool {
	return ch == "&" || ch == "|" || ch == "!" || ch == "NAND" || ch == "NOR" || ch == "XOR"
}

func precedence(op string) int {
	switch op {
	case "!":
		return 4
	case "NAND", "NOR", "XOR":
		return 3
	case "&":
		return 2
	case "|":
		return 1
	}
	return 0
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	for i := 0; i < len(input); i++ {
		ch := rune(input[i])
		if unicode.IsSpace(ch) {
			continue
		}
		if ch == '!' || ch == '&' || ch == '|' {
			tokens = append(tokens, Token{Type: OPERATOR, Value: string(ch)})
		} else if i+3 < len(input) && input[i:i+4] == "NAND" {
			tokens = append(tokens, Token{Type: OPERATOR, Value: "NAND"})
			i += 3 // Skip next three characters
		} else if i+2 < len(input) && input[i:i+3] == "NOR" {
			tokens = append(tokens, Token{Type: OPERATOR, Value: "NOR"})
			i += 2 // Skip next two characters
		} else if i+2 < len(input) && input[i:i+3] == "XOR" {
			tokens = append(tokens, Token{Type: OPERATOR, Value: "⊕"})
			i += 2 // Skip next two characters
		} else if unicode.IsLetter(ch) {
			tokens = append(tokens, Token{Type: VARIABLE, Value: string(ch)})
		} else {
			return nil, errors.New("invalid character in input")
		}
	}
	return tokens, nil
}

func shuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var operators []Token
	for _, token := range tokens {
		if token.Type == VARIABLE {
			output = append(output, token)
		} else if token.Type == OPERATOR {
			for len(operators) > 0 && precedence(operators[len(operators)-1].Value) >= precedence(token.Value) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}
	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output, nil
}

func buildParseTree(tokens []Token) (*Node, error) {
	var stack []*Node
	for _, token := range tokens {
		node := &Node{Value: token}
		if token.Type == OPERATOR {
			if len(stack) < 1 {
				return nil, errors.New("invalid expression")
			}
			if token.Value != "!" { // NOT is unary, others are binary
				if len(stack) < 2 {
					return nil, errors.New("invalid expression")
				}
				node.Right = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				node.Left = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else { // Handle unary NOT operator
				node.Left = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
		stack = append(stack, node)
	}
	if len(stack) != 1 {
		return nil, errors.New("invalid expression")
	}
	return stack[0], nil
}

func parse(input string) (*Node, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	postfixTokens, err := shuntingYard(tokens)
	if err != nil {
		return nil, err
	}
	return buildParseTree(postfixTokens)
}

func simplify(expr *Node) *Node {
	if expr == nil {
		return nil
	}

	// Recursively simplify the left and right subtrees
	expr.Left = simplify(expr.Left)
	expr.Right = simplify(expr.Right)

	// Apply De Morgan's laws to eliminate NAND and NOR
	switch expr.Value.Value {
	case "NAND":
		// Convert A NAND B to !(A&B)
		expr.Value.Value = "!"
		expr.Left = &Node{
			Value: Token{Type: OPERATOR, Value: "&"},
			Left:  expr.Left,
			Right: expr.Right,
		}
		expr.Right = nil
	case "NOR":
		// Convert A NOR B to !(A|B)
		expr.Value.Value = "!"
		expr.Left = &Node{
			Value: Token{Type: OPERATOR, Value: "|"},
			Left:  expr.Left,
			Right: expr.Right,
		}
		expr.Right = nil
	}

	return expr
}

func printExpr(expr *Node) string {
	exprf := ""
	if expr != nil {
		switch expr.Value.Type {
		case OPERATOR:
			if expr.Value.Value == "!" {
				// Unary operator
				return expr.Value.Value + printExpr(expr.Left)
			}
			// Binary operator
			return "(" + printExpr(expr.Left) + expr.Value.Value + printExpr(expr.Right) + ")"
		case VARIABLE:
			return expr.Value.Value
		default:
			return ""
		}
	}
	return exprf
}
