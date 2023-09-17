package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		handler    func(http.ResponseWriter, *http.Request)
		inData     string
		statusCode int
		outData    string
	}{
		{
			name:       "Initial POST",
			method:     http.MethodPost,
			target:     "/replace",
			handler:    replaceHandler,
			inData:     "Test data",
			statusCode: 200,
			outData:    "",
		},
		{
			name:       "GET after POST",
			method:     http.MethodGet,
			target:     "/get",
			handler:    getHandler,
			inData:     "",
			statusCode: 200,
			outData:    "Test data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, bytes.NewReader([]byte(tt.inData)))
			responseRecorder := httptest.NewRecorder()

			tt.handler(responseRecorder, request)

			if responseRecorder.Code != tt.statusCode {
				t.Errorf("Unexpected response code %d", responseRecorder.Code)
			}

			body, _ := io.ReadAll(responseRecorder.Body)
			textBody := string(body)

			if textBody != tt.outData {
				t.Errorf("Result mismatch: expected %s, got %s", tt.outData, textBody)
			}
		})
	}
}
