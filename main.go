package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/logic"
)

var cfAuthHeader string

func init() {
	cfAuthHeader = os.Getenv("CF_AUTH")
	if len(cfAuthHeader) == 0 {
		panic("CF_AUTH not set")
	}
}

func handleRequest(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	if request.Headers["x-tajnyzapis-cf-auth"] != cfAuthHeader {
		return common.MakeStringResponse("unauthorized", 401), nil
	}
	switch fmt.Sprintf("%s %s", request.RequestContext.HTTP.Method, request.RawPath) {

	case "POST /v1/zapis":
		return logic.CreateZapis(ctx, request)

	case "GET /v1/salt":
		return logic.GetSalt(ctx, request)

	case "GET /v1/zapis":
		return logic.GetZapis(ctx, request)

	}
	return common.MakeStringResponse("unknown method/path", 400), nil
}

func main() {
	lambda.Start(handleRequest)
}
