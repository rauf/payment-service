package validation

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type Errors struct {
	Errors []FieldError `json:"errors"`
}

func (e *Errors) IsValid() bool {
	return len(e.Errors) == 0
}

func (e *Errors) Add(field, message string) {
	e.Errors = append(e.Errors, FieldError{Field: field, Message: message})
}

func (e *Errors) Error() string {
	return "validation failed"
}
