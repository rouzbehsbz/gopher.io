package engineio

import "net/http"

type ErrorCode int

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

const (
	UnknownTransportErrorCode ErrorCode = iota
	UnknownSidErrorCode
	BadHandshakeMethodErrorCode
	BadRequestErrorCode
	ForbiddenErrorCode
	UnsupportedProtocolVersionErrorCode
)

var errorMessages = map[ErrorCode]string{
	UnknownTransportErrorCode:           "Transport unknown",
	UnknownSidErrorCode:                 "Session ID unknown",
	BadHandshakeMethodErrorCode:         "Bad handshake method",
	BadRequestErrorCode:                 "Bad request",
	ForbiddenErrorCode:                  "Forbidden",
	UnsupportedProtocolVersionErrorCode: "Unsupported protocol version",
}

var errorHttpStatusCodes = map[ErrorCode]int{
	UnknownTransportErrorCode:           http.StatusBadRequest,
	UnknownSidErrorCode:                 http.StatusBadRequest,
	BadHandshakeMethodErrorCode:         http.StatusBadRequest,
	BadRequestErrorCode:                 http.StatusBadRequest,
	ForbiddenErrorCode:                  http.StatusBadRequest,
	UnsupportedProtocolVersionErrorCode: http.StatusBadRequest,
}

func GetErrorMessage(errorCode ErrorCode) string {
	if msg, exists := errorMessages[errorCode]; exists {
		return msg
	}
	return "Unknown error"
}

func GetErrorHttpStatusCode(errorCode ErrorCode) int {
	if code, exists := errorHttpStatusCodes[errorCode]; exists {
		return code
	}
	return http.StatusInternalServerError
}
