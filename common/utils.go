package common

import (
	"math/rand"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func MakeStringResponse(what string, code int) *events.LambdaFunctionURLResponse {
	return &events.LambdaFunctionURLResponse{
		Body:       what,
		StatusCode: code,
	}
}
