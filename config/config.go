package config

const WordsInAccessKey = 3

const Base64Encoded32BytesLength = 44

// These are generated base base64 encoding 32 bytes of data.
const SaltLength = Base64Encoded32BytesLength
const AuthTokenLength = Base64Encoded32BytesLength
const HashedAccessKeyLength = Base64Encoded32BytesLength

const WrappedKeyLength = 64

const MaxFileSize = 10240 // 10 KB seems like a reasonable start point

const S3KeyLength = 36

const AwsRegion = "eu-north-1"
const BucketName = "tajny-zapis"
const TableName = "tajny-zapis"
