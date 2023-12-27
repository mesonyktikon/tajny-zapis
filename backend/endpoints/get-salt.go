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

	dynamoItems, _, err := storage.GetZapisOrDummyData(hashedAccessKey)
	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	var tollpass tokens.TollPass
	saltCandidates := make(map[string]string)
	tollpass.Candidates = make([]tokens.TollPassCandidate, len(dynamoItems))

	for idx, item := range dynamoItems {
		saltCandidates[item.SaltId] = item.Salt
		tollpass.Candidates[idx] = tokens.TollPassCandidate{
			AuthToken:  item.AuthToken,
			WrappedKey: item.WrappedKey,
			S3Key:      item.S3Key,
		}
	}

	tollpassJwt, err := tokens.GenerateTollPassJwt(&tollpass)

	if err != nil {
		logrus.Error(err)
		return common.MakeStringResponse("server error", 500), nil
	}

	return common.MakeJsonResponse(wire.GetSaltResponse{
		Salts:       saltCandidates,
		TollPassJwt: tollpassJwt,
	}, 200), nil
}
