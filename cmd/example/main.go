// Package main provides an example usage of the meridian package.
package main

import (
	"fmt"

	"github.com/matthalp/go-meridian"
)

func main() {
	fmt.Println("Meridian Package Example")
	fmt.Println("========================")
	fmt.Printf("Version: %s\n\n", meridian.Version)

	// Example 1: Using Greet
	greeting := meridian.Greet("Go Developer")
	fmt.Println(greeting)

	// Example 2: Using Example
	message := meridian.Example()
	fmt.Println(message)
}
