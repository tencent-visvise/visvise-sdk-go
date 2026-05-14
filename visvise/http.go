package visvise

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Environment represents the API environment
type Environment string

const (
	EnvProd Environment = "https://ws.visvise.com.cn"     // Production environment
	EnvTest Environment = "https://qa-ws.visvise.com.cn"  // Test environment
	EnvDev  Environment = "https://dev-ws.visvise.com.cn" // Development environment
)

// HTTPClient is the low-level HTTP client that handles signing and error handling
type HTTPClient struct {
	AppID     string
	SecretKey string
	UID       string
	BaseURL   string
	Timeout   int
	Client    *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(appID, secretKey, uid string, baseURL Environment, timeout int) *HTTPClient {
	return &HTTPClient{
		AppID:     appID,
		SecretKey: secretKey,
		UID:       uid,
		BaseURL:   strings.TrimRight(string(baseURL), "/"),
		Timeout:   timeout,
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// sign generates the HMAC-SHA256 signature
func (c *HTTPClient) sign(bodyStr string, ts string) string {
	signStr := bodyStr + ts
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

// buildHeaders builds the request headers with signature
func (c *HTTPClient) buildHeaders(bodyStr string) http.Header {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sign := c.sign(bodyStr, ts)

	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"app_id":       []string{c.AppID},
		"uid":          []string{c.UID},
		"ts":           []string{ts},
		"sign":         []string{sign},
	}
	return headers
}

// Post sends a POST request
func (c *HTTPClient) Post(path string, body interface{}) (interface{}, error) {
	var bodyStr string
	if body == nil {
		bodyStr = "{}"
	} else {
		// Serialize once to ensure consistency between signature and request
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyStr = string(jsonBytes)
	}

	url := c.BaseURL + "/" + strings.TrimLeft(path, "/")
	headers := c.buildHeaders(bodyStr)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, NewNetworkError(fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewNetworkError(fmt.Sprintf("failed to read response: %v", err))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, NewNetworkError(fmt.Sprintf("failed to parse response: %s", string(respBody)))
	}

	code, _ := result["code"].(float64)
	reqID, _ := result["req_id"].(string)
	msg, _ := result["msg"].(string)

	if code != 0 {
		return nil, RaiseForCode(int(code), msg, reqID)
	}

	return result["data"], nil
}

// SSEClient represents an SSE stream client
type SSEClient struct {
	HTTPClient  *HTTPClient
	Path        string
	Body        interface{}
	ReadTimeout int
}

// SSEEvent represents an SSE event
type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ReadTimeoutError represents SSE read timeout
type ReadTimeoutError struct {
	Message string
}

func (e *ReadTimeoutError) Error() string {
	return e.Message
}

// SSEIterator iterates over SSE events
type SSEIterator struct {
	client   *http.Client
	req      *http.Request
	resp     *http.Response
	reader   *bufio.Reader
	leftover string
}

// NewSSEIterator creates a new SSE iterator
func NewSSEIterator(httpClient *HTTPClient, path string, body interface{}, readTimeout int) (*SSEIterator, error) {
	var bodyStr string
	if body == nil {
		bodyStr = "{}"
	} else {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyStr = string(jsonBytes)
	}

	url := httpClient.BaseURL + "/" + strings.TrimLeft(path, "/")
	headers := httpClient.buildHeaders(bodyStr)
	headers.Set("Accept", "text/event-stream")

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	// Use a separate transport with connection timeout, but NO total Timeout.
	// http.Client.Timeout covers the entire request lifecycle including body read,
	// which is incompatible with SSE streaming. Instead we set TLSHandshakeTimeout
	// and use context for connection phase only.
	transport := &http.Transport{
		TLSHandshakeTimeout:   time.Duration(httpClient.Timeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(httpClient.Timeout) * time.Second,
	}
	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, NewNetworkError(fmt.Sprintf("SSE request failed: %v", err))
	}

	log.Printf("POST(SSE) response: status=%d, content-type=%s", resp.StatusCode, resp.Header.Get("Content-Type"))

	// Check HTTP status code like Python SDK does with raise_for_status()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		resp.Body.Close()
		bodyExcerpt := string(bodyBytes)
		if len(bodyExcerpt) > 500 {
			bodyExcerpt = bodyExcerpt[:500] + "...(truncated)"
		}
		return nil, NewNetworkError(fmt.Sprintf("SSE HTTP error [%d]: %s body=%s", resp.StatusCode, http.StatusText(resp.StatusCode), bodyExcerpt))
	}

	return &SSEIterator{
		client:   client,
		req:      req,
		resp:     resp,
		reader:   bufio.NewReader(resp.Body),
		leftover: "",
	}, nil
}

// Next returns the next SSE event
func (iter *SSEIterator) Next() (*SSEEvent, error) {
	event := ""
	var dataLines []string

	for {
		line, err := iter.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Process remaining data
				if len(dataLines) > 0 || event != "" {
					dataStr := strings.Join(dataLines, "\n")
					return &SSEEvent{Event: event, Data: parseSSELineData(dataStr)}, nil
				}
				return nil, io.EOF
			}
			return nil, NewNetworkError(fmt.Sprintf("SSE read error: %v", err))
		}

		lineStr := strings.TrimRight(line, "\r\n")

		if lineStr == "" {
			// Empty line indicates end of frame
			if event != "" || len(dataLines) > 0 {
				dataStr := strings.Join(dataLines, "\n")
				return &SSEEvent{Event: event, Data: parseSSELineData(dataStr)}, nil
			}
			continue
		}

		if strings.HasPrefix(lineStr, ":") {
			continue // Comment
		}

		if strings.HasPrefix(lineStr, "event:") {
			event = strings.TrimSpace(lineStr[6:])
		} else if strings.HasPrefix(lineStr, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(lineStr[5:]))
		}
	}
}

// Close closes the SSE stream
func (iter *SSEIterator) Close() error {
	if iter.resp != nil {
		return iter.resp.Body.Close()
	}
	return nil
}

func parseSSELineData(dataStr string) interface{} {
	var data interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return dataStr
	}
	return data
}

// PostSSE sends a POST request and returns an SSE iterator
func (c *HTTPClient) PostSSE(path string, body interface{}, readTimeout int) (*SSEIterator, error) {
	if readTimeout <= 0 {
		readTimeout = 1200 // Default 20 minutes
	}
	return NewSSEIterator(c, path, body, readTimeout)
}
