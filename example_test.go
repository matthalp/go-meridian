package meridian_test

import (
	"fmt"

	"github.com/matthalp/go-meridian"
)

func ExampleGreet() {
	message := meridian.Greet("World")
	fmt.Println(message)
	// Output: Hello, World!
}

func ExampleExample() {
	message := meridian.Example()
	fmt.Println(message)
	// Output: Welcome to Meridian!
}
