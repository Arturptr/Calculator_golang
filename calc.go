package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}
	return parseExpression(&tokens)
}

// Tokenizer: Converts the string expression into a list of tokens
func tokenize(expression string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder

	for _, ch := range expression {
		if unicode.IsDigit(ch) || ch == '.' {
			currentToken.WriteRune(ch) // Continue building number token
		} else if strings.ContainsRune("+-*/()", ch) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(ch)) // Add operator or bracket
		} else if unicode.IsSpace(ch) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else {
			return nil, errors.New("invalid character")
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens, nil
}

// Parse an expression (handles +, -, *, /, and parentheses)
func parseExpression(tokens *[]string) (float64, error) {
	return parseTerm(tokens)
}

// Parse a term (handles + and - operations)
func parseTerm(tokens *[]string) (float64, error) {
	result, err := parseFactor(tokens)
	if err != nil {
		return 0, err
	}

	for len(*tokens) > 0 {
		token := (*tokens)[0]
		if token != "+" && token != "-" {
			break
		}
		*tokens = (*tokens)[1:] // Consume the operator

		nextValue, err := parseFactor(tokens)
		if err != nil {
			return 0, err
		}

		if token == "+" {
			result += nextValue
		} else {
			result -= nextValue
		}
	}

	return result, nil
}

// Parse a factor (handles * and / operations)
func parseFactor(tokens *[]string) (float64, error) {
	result, err := parsePrimary(tokens)
	if err != nil {
		return 0, err
	}

	for len(*tokens) > 0 {
		token := (*tokens)[0]
		if token != "*" && token != "/" {
			break
		}
		*tokens = (*tokens)[1:] // Consume the operator

		nextValue, err := parsePrimary(tokens)
		if err != nil {
			return 0, err
		}

		if token == "*" {
			result *= nextValue
		} else {
			if nextValue == 0 {
				return 0, errors.New("division by zero")
			}
			result /= nextValue
		}
	}

	return result, nil
}

// Parse a primary value (number or parentheses)
func parsePrimary(tokens *[]string) (float64, error) {
	if len(*tokens) == 0 {
		return 0, errors.New("unexpected end of expression")
	}

	token := (*tokens)[0]
	*tokens = (*tokens)[1:] // Consume the token

	if token == "(" {
		result, err := parseTerm(tokens)
		if err != nil {
			return 0, err
		}
		if len(*tokens) == 0 || (*tokens)[0] != ")" {
			return 0, errors.New("missing closing parenthesis")
		}
		*tokens = (*tokens)[1:] // Consume closing parenthesis
		return result, nil
	}

	return strconv.ParseFloat(token, 64) // Parse a number
}
func main() {
	var expression string
	fmt.Print("Enter an expression: ")
	fmt.Scan(&expression)
	result, err := Calc(expression)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
