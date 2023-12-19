package endpoints

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mesonyktikon/tajny-zapis/common"
	"github.com/mesonyktikon/tajny-zapis/storage"
	"github.com/mesonyktikon/tajny-zapis/tokens"
	"github.com/mesonyktikon/tajny-zapis/wire"
	"github.com/sirupsen/logrus"
)

func GetZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	authToken := request.QueryStringParameters["authToken"]
	tollpassJwt := request.QueryStringParameters["tollpass"]

	tollpass, err := tokens.DecryptTollPassJwt(tollpassJwt)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse(err.Error(), 400), nil
	}

	if !tollpass.Valid || authToken != tollpass.AuthToken {
		logrus.Error(err)
		return common.MakeStringResponse("invalid", 400), nil
	}

	downloadUrl, err := storage.GeneratePresignedGetUrl(tollpass.S3Key)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	return common.MakeJsonResponse(wire.GetZapisResponse{
		DownloadUrl: downloadUrl,
		WrappedKey:  tollpass.WrappedKey,
	}, 200), nil
}
