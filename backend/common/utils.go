package common

import (
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		sb.WriteByte(charset[index])
	}
	return sb.String()
}

func MakeJsonResponse(what interface{}, code int) *events.LambdaFunctionURLResponse {
	body, err := json.Marshal(what)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal response body for error response")
		body = []byte("failed to marshal response body for error response")
		code = 500
	}

	return &events.LambdaFunctionURLResponse{
		Body:       string(body),
		StatusCode: code,
	}
}

func MakeStringResponse(what string, code int) *events.LambdaFunctionURLResponse {
	responseBody := struct {
		Msg string `json:"msg"`
	}{
		Msg: what,
	}

	body, err := json.Marshal(responseBody)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal response body for error response")
		body = []byte("failed to marshal response body for error response")
		code = 500
	}

	return &events.LambdaFunctionURLResponse{
		Body:       string(body),
		StatusCode: code,
	}
}
