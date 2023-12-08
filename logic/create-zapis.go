package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func CreateZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	var req common.CreateZapisRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return common.MakeStringResponse("malformed json", 400), nil
	}

	if len(req.Salt) != 36 {
		return common.MakeStringResponse("salt length != 36", 400), nil
	}

	if len(req.AuthToken) != 44 {
		return common.MakeStringResponse("authToken length != 44", 400), nil
	}

	if len(req.WrappedKey) != 64 {
		return common.MakeStringResponse("wrappedKey length != 64", 400), nil
	}

	if req.Ttl <= 0 {
		req.Ttl = 3600 * 24 * 7
	}

	req.Ttl = time.Now().Add(time.Duration(time.Second) * time.Duration(req.Ttl)).Unix()

	accessKey := common.GeneratePhrase(common.AccessKeyWords)
	s3Key := uuid.New().String()

	dynamoItem := common.TajnyZapisDynamoItem{
		Salt:       req.Salt,
		AccessKey:  accessKey,
		AuthToken:  req.AuthToken,
		WrappedKey: req.WrappedKey,
		S3Key:      s3Key,
		Ttl:        req.Ttl,
	}

	uploadUrl, err := storage.GeneratePresignedPutUrl(dynamoItem.S3Key)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	err = storage.MaybePutZapis(&dynamoItem)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	res := common.CreateZapisResponse{
		AccessKey: dynamoItem.AccessKey,
		UploadUrl: uploadUrl,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return common.MakeStringResponse("marshal error", 500), nil
	}

	return common.MakeStringResponse(string(resJson), 200), nil
}
