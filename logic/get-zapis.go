package logic

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func GetZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	authToken := request.QueryStringParameters["authToken"]
	tollpassJwt := request.QueryStringParameters["tollpass"]

	tollpass, err := DecryptTollPassJwt(tollpassJwt)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 400), nil
	}

	fmt.Printf("authToken: %s\ttollpass: %+v\n", authToken, tollpass)

	if !tollpass.Valid {
		return common.MakeStringResponse("invalid", 400), nil
	}

	if authToken != tollpass.AuthToken {
		return common.MakeStringResponse("bad auth token", 400), nil
	}

	downloadUrl, err := storage.GeneratePresignedGetUrl(tollpass.S3Key)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	return common.MakeJsonResponse(common.GetZapisResponse{
		DownloadUrl: downloadUrl,
		WrappedKey:  tollpass.WrappedKey,
	}, 200), nil
}
