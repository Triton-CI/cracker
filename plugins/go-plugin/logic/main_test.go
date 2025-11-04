package logic

import (
	"testing"
)

// FormatGreeting replicates the greeting logic from say_hello
func FormatGreeting(input string) string {
	return "👋🤗🎉 Extism is 💜 by " + input
}

// TestFormatGreeting tests the greeting formatting logic
func TestFormatGreeting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic input",
			input:    "World",
			expected: "👋🤗🎉 Extism is 💜 by World",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "👋🤗🎉 Extism is 💜 by ",
		},
		{
			name:     "special characters",
			input:    "Go Developers! 🚀",
			expected: "👋🤗🎉 Extism is 💜 by Go Developers! 🚀",
		},
		{
			name:     "long input",
			input:    "the amazing Go community around the world",
			expected: "👋🤗🎉 Extism is 💜 by the amazing Go community around the world",
		},
		{
			name:     "unicode input",
			input:    "世界",
			expected: "👋🤗🎉 Extism is 💜 by 世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FormatGreeting(tt.input)

			if output != tt.expected {
				t.Errorf("FormatGreeting(%q) = %q, want %q", tt.input, output, tt.expected)
			}
		})
	}
}
