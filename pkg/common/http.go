package common

import (
	apperrors "common/errors"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/url"
	"time"
)

func SerializeResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := toJSON(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       "{}", //TODO internal server error from string
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       jsonBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func SerializeError(err error) (events.APIGatewayProxyResponse, error) {
	jsonBody := ""
	statusCode := 500
	switch err.(type) {
	case *apperrors.Error:
		commonError := err.(*apperrors.Error)
		errorDto := ErrorResponseDto{
			ErrorCode:   commonError.ErrorCode,
			Description: commonError.Description,
			Params:      commonError.Params,
		}
		jsonBody, err = toJSON(errorDto)
		if err != nil {
			jsonBody = serializeInternalServerError()
		}
		statusCode = commonError.HttpStatusCode
		break
	default:
		jsonBody = serializeInternalServerError()
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       jsonBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func serializeInternalServerError() string {
	genericError := ErrorResponseDto{
		ErrorCode:   apperrors.INTERNAL_SERVER_ERROR,
		Description: "Internal server error has occurred",
	}
	jsonBody, err := toJSON(genericError)
	if err != nil {
		return fmt.Sprintf("{\"code\": \"%s\", \"description\": \"Internal server error has occurred\"}", apperrors.INTERNAL_SERVER_ERROR)
	}
	return jsonBody
}

func toJSON(input interface{}) (string, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type ErrorResponseDto struct {
	ErrorCode   string            `json:"errorCode"`
	Description string            `json:"description"`
	TraceId     string            `json:"traceId"`
	SpanId      string            `json:"spanId"`
	Params      map[string]string `json:"params,omitempty"`
}

func GetTimestampFromQueryParams(queryParams url.Values, param string) (*time.Time, error) {
	createdString := GetFilterByName(param, queryParams)
	var created *time.Time
	if len(createdString) > 0 {
		innerCreated, err := time.Parse(time.RFC3339, createdString)
		if err != nil {
			return nil, apperrors.InvalidRequestParameterWithValidation(fmt.Sprintf("Failed to parse query parameter %s with value '%s'", param, createdString), param, "RFC3339 timestamp", err)
		}
		created = &innerCreated
	}
	return created, nil
}

func GetFilterByName(name string, queryParams url.Values) string {
	filterArray := queryParams[name]
	filter := ""
	if filterArray != nil {
		filter = filterArray[0]
	}
	return filter
}
