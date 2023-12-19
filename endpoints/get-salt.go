package endpoints

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mesonyktikon/tajny-zapis/common"
	"github.com/mesonyktikon/tajny-zapis/storage"
	"github.com/mesonyktikon/tajny-zapis/tokens"
	"github.com/mesonyktikon/tajny-zapis/wire"
	"github.com/sirupsen/logrus"
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

	tollpassJwt, err := tokens.GenerateTollPassJwt(&tokens.TollPass{
		Valid:      valid,
		AuthToken:  dynamoItem.AuthToken,
		WrappedKey: dynamoItem.WrappedKey,
		S3Key:      dynamoItem.S3Key,
	})

	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	return common.MakeJsonResponse(wire.GetSaltResponse{
		Salt:        dynamoItem.Salt,
		TollPassJwt: tollpassJwt,
	}, 200), nil
}
