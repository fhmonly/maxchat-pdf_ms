package constants

type ErrorCode string

const (
	ERROR_INVALID_FILE_EXTENSION ErrorCode = "INVALID_FILE_EXTENSION"
	ERROR_FILE_TOO_LARGE         ErrorCode = "FILE_TOO_LARGE"
)
