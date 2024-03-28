package models

import "time"

type RequestLogData struct {
	Level     string `json:"level"`
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Header    string `json:"header"`
	Time      string `json:"time"`
	Host      string `json:"host"`
	Method    string `json:"method"`
	URI       string `json:"uri"`
	Body      string `json:"request_body,omitempty"`
	RemoteIP  string `json:"remote_ip"`
	BytesIn   int    `json:"bytes_in"`
}

type ResponseLogData struct {
	Level        string `json:"level"`
	Type         string `json:"type"`
	RequestID    string `json:"request_id"`
	Header       string `json:"header"`
	Time         string `json:"time"`
	URI          string `json:"uri"`
	Status       int    `json:"status"`
	Response     string `json:"response,omitempty"`
	Message      string `json:"message,omitempty"`
	LatencyHuman string `json:"latency_human,omitempty"`
	BytesOut     int    `json:"bytes_out,omitempty"`
}

type LogErrorCallStack struct {
	Level     string `json:"level"`
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Time      string `json:"time"`
	Message   string `json:"message"`
	URI       string `json:"uri"`
}

type RequestResponseBridge struct {
	RequestID string    `json:"request_id"`
	StartTime time.Time `json:"time"`
}
