package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"tuffbizz.com/m/v2/common"
	"tuffbizz.com/m/v2/storage"
)

func CreateZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	var req common.CreateZapisRequest
	err := json.Unmarshal([]byte(request.Body), &req)

	if err != nil {
		logrus.Error(err)
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
		return common.MakeStringResponse("fileSize <= 0", 400), nil
	}

	if req.FileSize > common.MaxFileSize {
		return common.MakeStringResponse(fmt.Sprintf("contentLength > %d", common.MaxFileSize), 400), nil
	}

	if req.Ttl <= 0 {
		req.Ttl = 3600 * 24 * 7
	}

	req.Ttl = time.Now().Add(time.Duration(time.Second) * time.Duration(req.Ttl)).Unix()

	zapis := &common.TajnyZapisDynamoItem{
		Salt:       req.Salt,
		AuthToken:  req.AuthToken,
		WrappedKey: req.WrappedKey,
		AccessKey:  common.GeneratePhrase(common.WordsInAccessKey),
		S3Key:      uuid.New().String(),
		Ttl:        req.Ttl,
	}

	uploadUrl, err := storage.GeneratePresignedPutUrl(zapis.S3Key, req.FileSize)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("failed to generate url", 500), nil
	}

	err = storage.MaybePutZapis(zapis)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("failed to put item", 500), nil
	}

	return common.MakeJsonResponse(common.CreateZapisResponse{
		AccessKey: zapis.AccessKey,
		UploadUrl: uploadUrl,
	}, 200), nil
}
