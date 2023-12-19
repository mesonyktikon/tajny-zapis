package logic

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func GetZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	authToken := request.QueryStringParameters["authToken"]
	tollpassJwt := request.QueryStringParameters["tollpass"]

	tollpass, err := DecryptTollPassJwt(tollpassJwt)
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

	return common.MakeJsonResponse(common.GetZapisResponse{
		DownloadUrl: downloadUrl,
		WrappedKey:  tollpass.WrappedKey,
	}, 200), nil
}
