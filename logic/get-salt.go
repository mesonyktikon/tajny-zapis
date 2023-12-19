package logic

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func GetSalt(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	accessKey := request.QueryStringParameters["accessKey"]

	if len(strings.Split(accessKey, " ")) < 3 {
		return common.MakeStringResponse("not enough words in access key", 400), nil
	}

	dynamoItem, valid, err := storage.GetZapisOrDummyData(accessKey)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	tollpassJwt, err := GenerateTollPassJwt(&TollPass{
		Valid:      valid,
		AuthToken:  dynamoItem.AuthToken,
		WrappedKey: dynamoItem.WrappedKey,
		S3Key:      dynamoItem.S3Key,
	})

	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	return common.MakeJsonResponse(common.GetSaltResponse{
		Salt:        dynamoItem.Salt,
		TollPassJwt: tollpassJwt,
	}, 200), nil
}
