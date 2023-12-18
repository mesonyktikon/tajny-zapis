package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func CreateZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	if request.IsBase64Encoded {
		return common.MakeStringResponse("base64 encoded request body not supported", 400), nil
	}

	var req common.CreateZapisRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return common.MakeStringResponse("malformed body", 400), nil
	}

	if len(req.Salt) != common.SaltLength {
		return common.MakeStringResponse(fmt.Sprintf("salt length != %d", common.SaltLength), 400), nil
	}

	if len(req.AuthToken) != common.AuthTokenLength {
		return common.MakeStringResponse(fmt.Sprintf("authToken length != %d", common.AuthTokenLength), 400), nil
	}

	if len(req.WrappedKey) != common.WrappedKeyLength {
		return common.MakeStringResponse(fmt.Sprintf("wrappedKey length != %d", common.WrappedKeyLength), 400), nil
	}

	if req.FileSize <= 0 {
		return common.MakeStringResponse("contentLength <= 0", 400), nil
	}

	if req.FileSize > common.MaxFileSize {
		return common.MakeStringResponse(fmt.Sprintf("contentLength > %d", common.MaxFileSize), 400), nil
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

	uploadUrl, err := storage.GeneratePresignedPutUrl(dynamoItem.S3Key, req.FileSize)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	err = storage.MaybePutZapis(&dynamoItem)
	if err != nil {
		return common.MakeStringResponse(err.Error(), 500), nil
	}

	return common.MakeJsonResponse(common.CreateZapisResponse{
		AccessKey: dynamoItem.AccessKey,
		UploadUrl: uploadUrl,
	}, 200), nil
}
