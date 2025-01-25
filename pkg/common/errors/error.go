package apperrors

import (
	"fmt"
	"net/http"
)

type Error struct {
	ErrorCode           string
	Description         string
	InternalDescription string
	Cause               error
	HttpStatusCode      int
	Params              map[string]string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, InternalDescription: %s, Cause: [%v]", e.ErrorCode, e.InternalDescription, e.Cause)
}

const (
	INTERNAL_SERVER_ERROR      = "INTERNAL_SERVER_ERROR"
	INVALID_REQUEST            = "INVALID_REQUEST"
	INVALID_REQUEST_PARAMETERS = "INVALID_REQUEST_PARAMETERS"
	ENTITY_NOT_FOUND           = "ENTITY_NOT_FOUND"
	ENTITY_ALREADY_EXIST       = "ENTITY_ALREADY_EXIST"
	INSUFFICIENT_PERMISSION    = "INSUFFICIENT_PERMISSION"
)

func Is(errorToCheck error, errorCode string) bool {
	if errorToCheck == nil {
		return false
	}
	convertedErrorToCheck, ok := errorToCheck.(*Error)
	if !ok {
		return false
	}
	return convertedErrorToCheck.ErrorCode == errorCode
}

func InternalServerError(message string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           INTERNAL_SERVER_ERROR,
		Description:         "Internal server error has occurred",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusInternalServerError,
	}
}

func InvalidRequest(message string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           INVALID_REQUEST,
		Description:         "Invalid request",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusBadRequest,
	}
}

func EntityNotFound(message string, key string, value string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           ENTITY_NOT_FOUND,
		Description:         "Entity not found",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusNotFound,
		Params: map[string]string{
			key: value,
		},
	}
}

func EntityNotFoundForMultipleFields(message string, params map[string]string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           ENTITY_NOT_FOUND,
		Description:         "Entity not found",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusNotFound,
		Params:              params,
	}
}

func InvalidRequestParameter(message string, paramName string) (error *Error) {
	return &Error{
		ErrorCode:           INVALID_REQUEST_PARAMETERS,
		Description:         "Invalid request parameter",
		InternalDescription: message,
		Cause:               nil,
		HttpStatusCode:      http.StatusBadRequest,
		Params: map[string]string{
			"param": paramName,
		},
	}
}

func InvalidRequestParameterWithValidation(message string, paramName string, rule string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           INVALID_REQUEST_PARAMETERS,
		Description:         "Invalid request parameter",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusBadRequest,
		Params: map[string]string{
			"param": paramName,
			"rule":  rule,
		},
	}
}

func EntityAlreadyExist(message string, key string, value string, cause error) (error *Error) {
	return &Error{
		ErrorCode:           ENTITY_ALREADY_EXIST,
		Description:         "Entity already exist",
		InternalDescription: message,
		Cause:               cause,
		HttpStatusCode:      http.StatusUnprocessableEntity,
		Params: map[string]string{
			key: value,
		},
	}
}

func UnauthorizedInsufficientPermissions(message string) (error *Error) {
	return &Error{
		ErrorCode:           INSUFFICIENT_PERMISSION,
		Description:         "You do not have rights to perform this action on this entity",
		InternalDescription: message,
		Cause:               nil,
		HttpStatusCode:      http.StatusForbidden,
	}
}
