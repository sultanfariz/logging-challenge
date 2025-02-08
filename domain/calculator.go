package domain

import "fmt"

// Calculator represents the calculator domain entity
type Calculator struct{}

// Calculate performs arithmetic operations
func (c *Calculator) Calculate(num1, num2 float64, op string) (float64, error) {
	switch op {
	case "add":
		return num1 + num2, nil
	case "sub":
		return num1 - num2, nil
	case "mult":
		return num1 * num2, nil
	case "div":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operation: %s", op)
	}
}
