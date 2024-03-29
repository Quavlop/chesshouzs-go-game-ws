package models

type Response struct {
	Status    int         `json:"status,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	ErrorData interface{} `json:"error_data,omitempty"`
}
