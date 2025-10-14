package utils

import (
	"regexp"
	"strings"
)

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail validates if a string is a valid email address
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Basic length check
	if len(email) < 5 || len(email) > 100 {
		return false
	}

	// Check for valid email format using regex
	return emailRegex.MatchString(email)
}

// IsValidEmailStrict validates email with additional checks
func IsValidEmailStrict(email string) bool {
	if !IsValidEmail(email) {
		return false
	}

	// Additional validations
	email = strings.TrimSpace(email)

	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return false
	}

	// Check for valid characters around @
	atIndex := strings.Index(email, "@")
	if atIndex <= 0 || atIndex >= len(email)-1 {
		return false
	}

	// Split local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]

	// Validate local part
	if len(local) == 0 || len(local) > 64 {
		return false
	}

	// Validate domain part
	if len(domain) == 0 || len(domain) > 255 {
		return false
	}

	// Domain should contain at least one dot
	if !strings.Contains(domain, ".") {
		return false
	}

	// Domain should not start or end with dot or hyphen
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
		return false
	}

	return true
}
