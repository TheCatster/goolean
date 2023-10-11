package goolean

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

			table, err := generateTruthTable(simplifiedTree)
			if err != nil {
				fmt.Println(err)
				continue
			}

			printTruthTable(table, simplifiedTree)
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
		return 4 // Complement variables
	case "NAND", "NOR", "XOR":
		return 3 // Combined logic gates
	case "&":
		return 2 // And precedence
	case "|":
		return 1 // Or has a precedence lower than other operators
	case "(":
		return 0 // Assigning a low precedence to open parenthesis
	case ")":
		return 0 // Assigning a low precedence to close parenthesis
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
			tokens = append(tokens, Token{Type: OPERATOR, Value: "XOR"})
			i += 2 // Skip next two characters
		} else if ch == '(' {
			tokens = append(tokens, Token{Type: OPERATOR, Value: "("})
		} else if ch == ')' {
			tokens = append(tokens, Token{Type: OPERATOR, Value: ")"})
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
			if token.Value == "(" {
				operators = append(operators, token)
			} else if token.Value == ")" {
				for len(operators) > 0 && operators[len(operators)-1].Value != "(" {
					output = append(output, operators[len(operators)-1])
					operators = operators[:len(operators)-1]
				}
				if len(operators) == 0 {
					return nil, errors.New("mismatched parentheses")
				}
				operators = operators[:len(operators)-1] // Pop the open parenthesis
			} else {
				for len(operators) > 0 && precedence(operators[len(operators)-1].Value) >= precedence(token.Value) {
					output = append(output, operators[len(operators)-1])
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, token)
			}
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

func evaluateExpression(expr *Node, variables []string, values []bool) (bool, error) {
	if expr == nil {
		return false, errors.New("expression cannot be nil")
	}

	if expr.Value.Type == VARIABLE {
		// Find index of variable in variables slice
		index := indexOf(variables, expr.Value.Value)
		if index == -1 {
			return false, errors.New("variable not found: " + expr.Value.Value)
		}
		return values[index], nil
	}

	if expr.Value.Type == OPERATOR {
		switch expr.Value.Value {
		case "!":
			value, err := evaluateExpression(expr.Left, variables, values)
			if err != nil {
				return false, err
			}
			return !value, nil
		case "&":
			leftValue, err := evaluateExpression(expr.Left, variables, values)
			if err != nil {
				return false, err
			}
			rightValue, err := evaluateExpression(expr.Right, variables, values)
			if err != nil {
				return false, err
			}
			return leftValue && rightValue, nil
		case "XOR":
			leftValue, err := evaluateExpression(expr.Left, variables, values)
			if err != nil {
				return false, err
			}
			rightValue, err := evaluateExpression(expr.Right, variables, values)
			if err != nil {
				return false, err
			}
			return leftValue != rightValue, nil
		}
	}

	return false, errors.New("invalid node type")
}

func getUniqueVariables(expr *Node) []string {
	if expr == nil {
		return []string{}
	}

	var variables []string

	if expr.Value.Type == VARIABLE {
		variables = append(variables, expr.Value.Value)
	}

	leftVariables := getUniqueVariables(expr.Left)
	rightVariables := getUniqueVariables(expr.Right)

	// Combine and deduplicate variables
	variables = append(variables, leftVariables...)
	variables = append(variables, rightVariables...)
	variables = deduplicate(variables)

	return variables
}

func deduplicate(items []string) []string {
	seen := map[string]bool{}
	unique := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			unique = append(unique, item)
		}
	}

	return unique
}

func indexOf(slice []string, item string) int {
	for i, value := range slice {
		if value == item {
			return i
		}
	}
	return -1
}

func generateTruthTable(expr *Node) ([][]bool, error) {
	// Get unique variables from expression
	variables := getUniqueVariables(expr)
	numCombinations := 1 << len(variables)
	truthTable := make([][]bool, numCombinations)

	for i := 0; i < numCombinations; i++ {
		row := make([]bool, len(variables)+1) // +1 for the result column
		for j, _ := range variables {
			row[j] = (i & (1 << j)) > 0
		}
		value, err := evaluateExpression(expr, variables, row)
		if err != nil {
			return nil, err
		}
		row[len(variables)] = value
		truthTable[i] = row
	}

	return truthTable, nil
}

func printTruthTable(truthTable [][]bool, expr *Node) {
	variables := getUniqueVariables(expr)
	headers := append(variables, printExpr(expr))
	fmt.Println(strings.Join(headers, " | "))

	for _, row := range truthTable {
		strRow := make([]string, len(row))
		for i, value := range row {
			strRow[i] = fmt.Sprintf("%v", value)
		}
		fmt.Println(strings.Join(strRow, " | "))
	}
}
