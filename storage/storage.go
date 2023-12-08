package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"tuffbizz.com/m/v2/common"
)

var sess *session.Session
var ddb *dynamodb.DynamoDB
var s3Client *s3.S3

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ddb = dynamodb.New(sess)
	s3Client = s3.New(sess)
}

func stringAttribute(str string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: aws.String(str),
	}
}

func numberAttribute(num int64) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(num, 10)),
	}
}

func MaybePutZapis(zapis *common.TajnyZapisDynamoItem) error {
	conditionExpression := "attribute_not_exists(salt)"

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"salt":        stringAttribute(zapis.Salt),
			"access_key":  stringAttribute(zapis.AccessKey),
			"auth_token":  stringAttribute(zapis.AuthToken),
			"wrapped_key": stringAttribute(zapis.WrappedKey),
			"s3_key":      stringAttribute(zapis.S3Key),
			"ttl":         numberAttribute(zapis.Ttl),
		},
		TableName:           aws.String("tajny-zapis"),
		ConditionExpression: aws.String(conditionExpression),
	}

	_, err := ddb.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// Returns dummy data if an item is not found.
func MaybeGetZapis(accessKey string) (*common.TajnyZapisDynamoItem, error) {
	items, err := FetchFromGSI("access_key-index", "access_key", accessKey)
	if err != nil {
		return nil, err
	}

	if len(items) > 1 {
		panic(fmt.Sprintf("access_key-index has more than one item for access_key=%s", accessKey))
	}

	if len(items) == 0 {
		return &common.TajnyZapisDynamoItem{
			Salt:       common.GenerateRandomString(common.SaltLength),
			AccessKey:  accessKey,
			AuthToken:  common.GenerateRandomString(common.AuthTokenLength),
			WrappedKey: common.GenerateRandomString(common.WrappedKeyLength),
			S3Key:      common.GenerateRandomString(common.S3KeyLength),
			Ttl:        time.Now().Unix(),
		}, nil
	}

	return items[0], nil
}

func FetchFromGSI(indexName, partitionKey, partitionValue string) ([]*common.TajnyZapisDynamoItem, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("tajny-zapis"),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("#partitionKey = :partitionValue"),
		ExpressionAttributeNames: map[string]*string{
			"#partitionKey": aws.String(partitionKey),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionValue": stringAttribute(partitionValue),
		},
	}

	result, err := ddb.Query(queryInput)
	if err != nil {
		return nil, err
	}

	var items []*common.TajnyZapisDynamoItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func GeneratePresignedPutUrl(s3Key string) (string, error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("tajny-zapis"),
		Key:    aws.String(s3Key),
	})

	url, err := req.Presign(1 * time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
