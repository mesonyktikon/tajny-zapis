package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/logic"
)

func doTest() string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ddb := dynamodb.New(sess)
	conditionExpression := "attribute_not_exists(accessKey)"

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"accessKey": {
				S: aws.String("test"),
			},
		},
		TableName:           aws.String("tajny-zapis"),
		ConditionExpression: aws.String(conditionExpression),
	}

	_, err := ddb.PutItem(input)
	if err != nil {
		return err.Error()
	}

	return "ok"
}
func handleRequest(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	if request.Headers["x-tajnyzapis-cf-auth"] != "1502b061-b4e5-46fa-ae8b-7cfe562d41b1" {
		return common.MakeStringResponse("unauthorized", 401), nil
	}

	switch fmt.Sprintf("%s %s", request.RequestContext.HTTP.Method, request.RawPath) {
	case "GET /v1/test":
		return common.MakeStringResponse(doTest(), 200), nil
	case "GET /v1/salt":
		return logic.GetSalt(ctx, request)
	case "POST /v1/zapis":
		return logic.CreateZapis(ctx, request)
	}

	return common.MakeStringResponse("unknown method/path", 400), nil
}

func main() {
	lambda.Start(handleRequest)
}
