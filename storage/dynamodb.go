package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mesonyktikon/tajny-zapis/common"
	"github.com/mesonyktikon/tajny-zapis/config"
)

type TajnyZapisDynamoItem struct {
	Salt      string `dynamodbav:"salt"`
	AccessKey string `dynamodbav:"access_key"`

	AuthToken  string `dynamodbav:"auth_token"`
	WrappedKey string `dynamodbav:"wrapped_key"`

	S3Key string `dynamodbav:"s3_key"`
	Ttl   int64  `dynamodbav:"ttl"`
}

var ddb *dynamodb.DynamoDB

func init() {
	ddb = dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(config.AwsRegion),
		},
	})))
}

func TryPutZapis(zapis *TajnyZapisDynamoItem) error {
	item, err := dynamodbattribute.MarshalMap(zapis)
	if err != nil {
		return err
	}

	_, err = ddb.PutItem(&dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(config.TableName),
		ConditionExpression: aws.String("attribute_not_exists(salt)"),
	})
	return err
}

func GetZapisOrDummyData(accessKey string) (*TajnyZapisDynamoItem, bool, error) {
	items, err := fetchFromGSI("access_key-index", "access_key", accessKey)
	if err != nil {
		return nil, false, err
	}

	if len(items) > 1 {
		panic(fmt.Sprintf("access_key-index has more than one item for access_key=%s", accessKey))
	}

	if len(items) == 0 {
		return &TajnyZapisDynamoItem{
			Salt:       common.GenerateRandomString(config.SaltLength),
			AccessKey:  accessKey,
			AuthToken:  common.GenerateRandomString(config.AuthTokenLength),
			WrappedKey: common.GenerateRandomString(config.WrappedKeyLength),
			S3Key:      common.GenerateRandomString(config.S3KeyLength),
			Ttl:        time.Now().Unix(),
		}, false, nil
	}

	return items[0], true, nil
}

func fetchFromGSI(indexName, partitionKey, partitionValue string) ([]*TajnyZapisDynamoItem, error) {
	result, err := ddb.Query(&dynamodb.QueryInput{
		TableName:              aws.String(config.TableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("#partitionKey = :partitionValue"),
		ExpressionAttributeNames: map[string]*string{
			"#partitionKey": aws.String(partitionKey),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionValue": stringAttribute(partitionValue),
		},
	})
	if err != nil {
		return nil, err
	}

	var items []*TajnyZapisDynamoItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	return items, err
}

func stringAttribute(str string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: aws.String(str),
	}
}
