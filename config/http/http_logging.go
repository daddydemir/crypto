package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	next http.RoundTripper
}

func NewLoggingMiddleware(next http.RoundTripper) *LoggingMiddleware {
	if next == nil {
		next = http.DefaultTransport
	}
	return &LoggingMiddleware{
		next: next,
	}
}

func (m *LoggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := m.next
	if transport == nil {
		transport = http.DefaultTransport
	}

	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}

	log.Printf("=== HTTP REQUEST ===")
	log.Printf("Time: %s", time.Now().Format(time.RFC3339))
	log.Printf("Method: %s", req.Method)
	log.Printf("URL: %s", req.URL.String())
	log.Printf("Headers: %v", req.Header)

	if len(reqBodyBytes) > 0 {
		log.Printf("Request Body:")
		if isJSON(reqBodyBytes) {
			log.Printf("%s", prettyJSON(reqBodyBytes))
		} else {
			log.Printf("%s", string(reqBodyBytes))
		}
	} else {
		log.Printf("Request Body: [EMPTY]")
	}

	start := time.Now()
	resp, err := transport.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("=== REQUEST ERROR ===")
		log.Printf("Error: %v", err)
		log.Printf("Duration: %v", duration)
		return nil, err
	}

	var respBodyBytes []byte
	if resp.Body != nil {
		respBodyBytes, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))
	}

	log.Printf("=== HTTP RESPONSE ===")
	log.Printf("Time: %s", time.Now().Format(time.RFC3339))
	log.Printf("Status: %d %s", resp.StatusCode, resp.Status)
	log.Printf("Duration: %v", duration)
	log.Printf("Response Headers: %v", resp.Header)

	if len(respBodyBytes) > 0 {
		log.Printf("Response Body:")
		if isJSON(respBodyBytes) {
			log.Printf("%s", prettyJSON(respBodyBytes))
		} else {
			log.Printf("%s", string(respBodyBytes))
		}
	} else {
		log.Printf("Response Body: [EMPTY]")
	}

	log.Printf("=== END ===\n\n")

	return resp, nil

}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

func prettyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data)
	}
	return prettyJSON.String()
}
