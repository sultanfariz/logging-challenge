package domain

import "fmt"

// Greeting represents the greeting domain entity
type Greeting struct{}

// GetGreeting returns a greeting message
func (g *Greeting) GetGreeting(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}
	return fmt.Sprintf("Hello, %s!", name), nil
}
