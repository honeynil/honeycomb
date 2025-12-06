package main

import "fmt"

func FormatResult(operation string, a, b, result float64) string {
	var symbol string
	switch operation {
	case "add":
		symbol = "+"
	case "subtract":
		symbol = "-"
	case "multiply":
		symbol = "×"
	case "divide":
		symbol = "÷"
	default:
		symbol = "?"
	}

	return fmt.Sprintf("%.2f %s %.2f = %.2f", a, symbol, b, result)
}
