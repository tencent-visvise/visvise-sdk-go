package visvise

import "fmt"

// WeaverError is the base error type for all SDK errors
type WeaverError struct {
	Code    int    // Error code
	ReqID   string // Request ID
	Message string // Error message
}

func (e *WeaverError) Error() string {
	if e.ReqID != "" {
		return fmt.Sprintf("[%d] %s (req_id=%s)", e.Code, e.Message, e.ReqID)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NetworkError indicates network failure (connection timeout, DNS resolution failed, etc.)
type NetworkError struct {
	WeaverError
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error: %s", e.Message)
}

// SignatureError indicates signature mismatch (HTTP 400)
type SignatureError struct {
	WeaverError
}

// InvalidParamsError indicates invalid request parameters (120008)
type InvalidParamsError struct {
	WeaverError
}

// UserNotFoundError indicates user not found (120017)
type UserNotFoundError struct {
	WeaverError
}

// PermissionDeniedError indicates permission denied (120018)
type PermissionDeniedError struct {
	WeaverError
}

// QuotaExceededError indicates daily generation quota exceeded (120020)
type QuotaExceededError struct {
	WeaverError
}

// ProjectPermissionError indicates project permission not authorized (120027)
type ProjectPermissionError struct {
	WeaverError
}

// ServerNetworkError indicates server network error (120028)
type ServerNetworkError struct {
	WeaverError
}

// ServerTimeoutError indicates server processing timeout (120032)
type ServerTimeoutError struct {
	WeaverError
}

// RateLimitError indicates too many requests (120040)
type RateLimitError struct {
	WeaverError
}

// ModelGenerationError indicates model generation failed (async task status=4)
type ModelGenerationError struct {
	WeaverError
	ModelID string
}

func (e *ModelGenerationError) Error() string {
	if e.ModelID != "" {
		return fmt.Sprintf("[%d] %s (model_id=%s)", e.Code, e.Message, e.ModelID)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// PollingTimeoutError indicates waiting for model completion timeout
type PollingTimeoutError struct {
	WeaverError
	ModelID string
	Timeout int
}

func (e *PollingTimeoutError) Error() string {
	return fmt.Sprintf("Timeout waiting for model %s (timeout=%ds)", e.ModelID, e.Timeout)
}

// errorCodeMapping maps error codes to error types
var errorCodeMapping = map[int]func(WeaverError) error{
	400:    func(e WeaverError) error { return &SignatureError{WeaverError: e} },
	120008: func(e WeaverError) error { return &InvalidParamsError{WeaverError: e} },
	120017: func(e WeaverError) error { return &UserNotFoundError{WeaverError: e} },
	120018: func(e WeaverError) error { return &PermissionDeniedError{WeaverError: e} },
	120020: func(e WeaverError) error { return &QuotaExceededError{WeaverError: e} },
	120027: func(e WeaverError) error { return &ProjectPermissionError{WeaverError: e} },
	120028: func(e WeaverError) error { return &ServerNetworkError{WeaverError: e} },
	120032: func(e WeaverError) error { return &ServerTimeoutError{WeaverError: e} },
	120040: func(e WeaverError) error { return &RateLimitError{WeaverError: e} },
}

// RaiseForCode creates and returns the appropriate error based on the error code
func RaiseForCode(code int, msg string, reqID string) error {
	weaverErr := WeaverError{Code: code, Message: msg, ReqID: reqID}
	if constructor, ok := errorCodeMapping[code]; ok {
		return constructor(weaverErr)
	}
	return &weaverErr
}

// NewWeaverError creates a new WeaverError
func NewWeaverError(code int, message string, reqID string) *WeaverError {
	return &WeaverError{Code: code, Message: message, ReqID: reqID}
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(message string) *NetworkError {
	return &NetworkError{WeaverError: WeaverError{Code: -1, Message: message}}
}

// NewModelGenerationError creates a new ModelGenerationError
func NewModelGenerationError(message string, code int, modelID string, reqID string) *ModelGenerationError {
	return &ModelGenerationError{
		WeaverError: WeaverError{Code: code, Message: message, ReqID: reqID},
		ModelID:     modelID,
	}
}

// NewPollingTimeoutError creates a new PollingTimeoutError
func NewPollingTimeoutError(modelID string, timeout int) *PollingTimeoutError {
	return &PollingTimeoutError{
		WeaverError: WeaverError{Code: -2, Message: fmt.Sprintf("Timeout waiting for model %s", modelID)},
		ModelID:     modelID,
		Timeout:     timeout,
	}
}
