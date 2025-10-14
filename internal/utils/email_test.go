package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Valid email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "Valid email with numbers",
			email:    "user123@example123.com",
			expected: true,
		},
		{
			name:     "Valid email with special characters",
			email:    "user.name+tag@example.com",
			expected: true,
		},
		{
			name:     "Invalid email - no @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "Invalid email - no domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "Invalid email - no local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "Invalid email - no TLD",
			email:    "test@example",
			expected: false,
		},
		{
			name:     "Invalid email - multiple @",
			email:    "test@@example.com",
			expected: false,
		},
		{
			name:     "Empty email",
			email:    "",
			expected: false,
		},
		{
			name:     "Email with spaces",
			email:    " test@example.com ",
			expected: true, // Should be trimmed
		},
		{
			name:     "Invalid email - too long",
			email:    "verylongusernamethatexceedsthelimitverylongusernamethatexceedsthelimitverylongusernamethatexceedsthelimitverylongusernamethatexceedsthelimit@example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			assert.Equal(t, tt.expected, result, "Expected %v for email: %s", tt.expected, tt.email)
		})
	}
}

func TestIsValidEmailStrict(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Invalid email - consecutive dots",
			email:    "test..name@example.com",
			expected: false,
		},
		{
			name:     "Invalid email - domain starts with dot",
			email:    "test@.example.com",
			expected: false,
		},
		{
			name:     "Invalid email - domain ends with dot",
			email:    "test@example.com.",
			expected: false,
		},
		{
			name:     "Invalid email - domain starts with hyphen",
			email:    "test@-example.com",
			expected: false,
		},
		{
			name:     "Invalid email - domain ends with hyphen",
			email:    "test@example.com-",
			expected: false,
		},
		{
			name:     "Invalid email - no dot in domain",
			email:    "test@example",
			expected: false,
		},
		{
			name:     "Invalid email - local part too long",
			email:    "verylongusernamethatexceedsthelimitverylongusernamethatexceedsthelimit@example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmailStrict(tt.email)
			assert.Equal(t, tt.expected, result, "Expected %v for email: %s", tt.expected, tt.email)
		})
	}
}
