package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func GetSalt(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	accessKey := request.QueryStringParameters["accessKey"]

	if len(strings.Split(accessKey, " ")) < 3 {
		return common.MakeStringResponse("not enough words in access key", 400), nil
	}

	dynamoItem, err := storage.MaybeGetZapis(accessKey)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	tollpass := common.TollPass{
		Valid:      accessKey == dynamoItem.AccessKey,
		Salt:       dynamoItem.Salt,
		WrappedKey: dynamoItem.WrappedKey,
		AuthToken:  dynamoItem.AuthToken,
		S3Key:      dynamoItem.S3Key,
	}

	tollpassJwt, err := GenerateTollPassJwt(&tollpass)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	res := common.GetSaltResponse{
		Salt:        dynamoItem.Salt,
		TollPassJwt: tollpassJwt,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return common.MakeStringResponse("marshal error", 500), nil
	}

	return common.MakeStringResponse(string(resJson), 200), nil
}
