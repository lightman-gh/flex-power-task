package response

import (
	"fmt"
	"net/http"
)

type Fail struct {
	Message    string `json:"message"`
	StatusCode uint32 `json:"status_code"`
}

type Success struct {
	any
}

func OK(data any) *Success {
	return &Success{
		data,
	}
}

func Error(format string, args ...any) *Fail {
	return &Fail{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf(format, args...),
	}
}
