package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mesonyktikon/tajny-zapis/common"
	"github.com/mesonyktikon/tajny-zapis/config"
	"github.com/mesonyktikon/tajny-zapis/storage"
	"github.com/mesonyktikon/tajny-zapis/wire"
	"github.com/sirupsen/logrus"
)

func CreateZapis(ctx context.Context, request *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	var req wire.CreateZapisRequest
	err := json.Unmarshal([]byte(request.Body), &req)

	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("malformed body", 400), nil
	}

	if len(req.Salt) != config.SaltLength {
		return common.MakeStringResponse(fmt.Sprintf("salt length != %d", config.SaltLength), 400), nil
	}

	if len(req.AuthToken) != config.AuthTokenLength {
		return common.MakeStringResponse(fmt.Sprintf("authToken length != %d", config.AuthTokenLength), 400), nil
	}

	if len(req.WrappedKey) != config.WrappedKeyLength {
		return common.MakeStringResponse(fmt.Sprintf("wrappedKey length != %d", config.WrappedKeyLength), 400), nil
	}

	if req.FileSize <= 0 {
		return common.MakeStringResponse("fileSize <= 0", 400), nil
	}

	if req.FileSize > config.MaxFileSize {
		return common.MakeStringResponse(fmt.Sprintf("contentLength > %d", config.MaxFileSize), 400), nil
	}

	if req.Ttl <= 0 {
		req.Ttl = 3600 * 24 * 7
	}

	req.Ttl = time.Now().Add(time.Duration(time.Second) * time.Duration(req.Ttl)).Unix()

	zapis := &storage.TajnyZapisDynamoItem{
		Salt:       req.Salt,
		AuthToken:  req.AuthToken,
		WrappedKey: req.WrappedKey,
		AccessKey:  common.GeneratePhrase(config.WordsInAccessKey),
		S3Key:      common.GenerateRandomString(config.S3KeyLength),
		Ttl:        req.Ttl,
	}

	uploadUrl, err := storage.GeneratePresignedPutUrl(zapis.S3Key, req.FileSize)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("failed to generate url", 500), nil
	}

	err = storage.TryPutZapis(zapis)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("failed to put item", 500), nil
	}

	return common.MakeJsonResponse(wire.CreateZapisResponse{
		AccessKey: zapis.AccessKey,
		UploadUrl: uploadUrl,
	}, 200), nil
}
