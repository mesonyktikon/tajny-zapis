package endpoints

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mesonyktikon/tajny-zapis/common"
	"github.com/mesonyktikon/tajny-zapis/config"
	"github.com/mesonyktikon/tajny-zapis/storage"
	"github.com/mesonyktikon/tajny-zapis/tokens"
	"github.com/mesonyktikon/tajny-zapis/wire"
	"github.com/sirupsen/logrus"
)

func GetSalt(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	hashedAccessKey := request.Headers["x-tajny-zapis-hashed-access-key"]

	if len(hashedAccessKey) != config.HashedAccessKeyLength {
		return common.MakeStringResponse("incorrect length for hashed access key", 400), nil
	}

	dynamoItem, valid, err := storage.GetZapisOrDummyData(hashedAccessKey)
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
