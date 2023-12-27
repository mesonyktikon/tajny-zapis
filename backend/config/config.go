package config

const (
	WordsInAccessKey           = 3
	Base64Encoded32BytesLength = 44

	// These are generated base base64 encoding 32 bytes of data.
	SaltLength            = Base64Encoded32BytesLength
	AuthTokenLength       = Base64Encoded32BytesLength
	HashedAccessKeyLength = Base64Encoded32BytesLength
	SaltIdLength          = Base64Encoded32BytesLength

	WrappedKeyLength = 64

	MaxFileSize = 10240 // 10 KB seems like a reasonable start point

	S3KeyLength = 36

	BucketRegion = "eu-north-1"
	BucketName   = "tajny-zapis"

	TableRegion = "eu-north-1"
	TableName   = "tajny-zapis"
)
