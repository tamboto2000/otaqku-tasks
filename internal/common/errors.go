package common

type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

type FieldError struct {
	Name     string   `json:"name"`
	Messages []string `json:"messages"`
}

// Common error codes
const (
	ErrCodeNotFound        = "not_foud"
	ErrCodeInputValidation = "input_validation"
	ErrCodeAlreadyExists   = "already_exists"
	ErrCodeUnauthorized    = "unauthorized"
)

// Common errors
var (
	ErrNotFound = Error{
		Code:    ErrCodeNotFound,
		Message: "Not found",
	}
)
