package errors

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Padding(1)
)

type AppError struct {
	Type    string
	Message string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewNetworkError(message string) AppError {
	return AppError{
		Type:    "Network Error",
		Message: message,
	}
}

func NewAuthenticationError(message string) AppError {
	return AppError{
		Type:    "Authentication Error",
		Message: message,
	}
}

func NewInputError(message string) AppError {
	return AppError{
		Type:    "Input Error",
		Message: message,
	}
}

func NewAPIError(message string) AppError {
	return AppError{
		Type:    "API Error",
		Message: message,
	}
}

func FormatError(err error) string {
	if appErr, ok := err.(AppError); ok {
		return ErrorStyle.Render(fmt.Sprintf("%s\n%s", appErr.Type, appErr.Message))
	}
	return ErrorStyle.Render(err.Error())
}
