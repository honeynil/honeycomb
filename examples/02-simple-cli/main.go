package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: add requires two numbers")
			fmt.Fprintln(os.Stderr, "Usage: calc add <number1> <number2>")
			os.Exit(1)
		}
		a, b, err := parseNumbers(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		result := Add(a, b)
		fmt.Println(FormatResult("add", a, b, result))

	case "subtract":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: subtract requires two numbers")
			os.Exit(1)
		}
		a, b, err := parseNumbers(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		result := Subtract(a, b)
		fmt.Println(FormatResult("subtract", a, b, result))

	case "multiply":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: multiply requires two numbers")
			os.Exit(1)
		}
		a, b, err := parseNumbers(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		result := Multiply(a, b)
		fmt.Println(FormatResult("multiply", a, b, result))

	case "divide":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: divide requires two numbers")
			os.Exit(1)
		}
		a, b, err := parseNumbers(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		result, err := Divide(a, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(FormatResult("divide", a, b, result))

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(2)
	}
}

func parseNumbers(a, b string) (float64, float64, error) {
	num1, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first number: %s", a)
	}
	num2, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second number: %s", b)
	}
	return num1, num2, nil
}

func printHelp() {
	fmt.Println("Simple Calculator CLI")
	fmt.Println()
	fmt.Println("Usage: calc <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <a> <b>        Add two numbers")
	fmt.Println("  subtract <a> <b>   Subtract b from a")
	fmt.Println("  multiply <a> <b>   Multiply two numbers")
	fmt.Println("  divide <a> <b>     Divide a by b")
	fmt.Println("  help               Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  calc add 5 3")
	fmt.Println("  calc multiply 4.5 2")
	fmt.Println("  calc divide 10 3")
}
