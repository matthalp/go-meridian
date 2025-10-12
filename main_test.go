package meridian

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple greeting",
			input:    "World",
			expected: "Hello, World!",
		},
		{
			name:     "greeting with name",
			input:    "Alice",
			expected: "Hello, Alice!",
		},
		{
			name:     "empty name",
			input:    "",
			expected: "Hello, !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greet(tt.input)
			if result != tt.expected {
				t.Errorf("Greet(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExample(t *testing.T) {
	result := Example()
	expected := "Welcome to Meridian!"
	if result != expected {
		t.Errorf("Example() = %q, want %q", result, expected)
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}
