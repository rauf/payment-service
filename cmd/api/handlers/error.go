package handlers

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Err     error  `json:"-"`
}

func NewResponse(code int, message string, data any, err error) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
		Err:     err,
	}
}
