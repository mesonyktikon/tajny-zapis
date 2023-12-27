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
		return common.MakeStringResponse("invalid tollpass", 401), nil
	}

	for _, candidate := range tollpass.Candidates {
		if authToken == candidate.AuthToken {
			downloadUrl, err := storage.GeneratePresignedGetUrl(candidate.S3Key)
			if err != nil {
				logrus.Error(err)
				return common.MakeStringResponse("server error", 500), nil
			}

			return common.MakeJsonResponse(wire.GetZapisResponse{
				DownloadUrl: downloadUrl,
				WrappedKey:  candidate.WrappedKey,
			}, 200), nil
		}
	}

	return common.MakeStringResponse("invalid salt or auth token", 403), nil
}
