package main

import (
	"context"

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

func handleGet(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	switch request.RawPath {
	case "/":
		return common.MakeStringResponse(doTest(), 200), nil
	case "/salt":
		return logic.GetSalt(ctx, request)
	}
	return common.MakeStringResponse("unknown path", 400), nil
}

func handlePost(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	switch request.RawPath {
	case "/zapis":
		return logic.CreateZapis(ctx, request)
	}
	return common.MakeStringResponse("unknown path", 400), nil
}

func handleRequest(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	switch request.RequestContext.HTTP.Method {
	case "GET":
		return handleGet(ctx, request)
	case "POST":
		return handlePost(ctx, request)
	}

	return common.MakeStringResponse("unknown method", 400), nil
}

func main() {
	lambda.Start(handleRequest)
}
