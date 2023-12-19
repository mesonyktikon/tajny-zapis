package common

type CreateZapisRequest struct {
	Salt       string `json:"salt"`
	AuthToken  string `json:"authToken"`
	WrappedKey string `json:"wrappedKey"`
	FileSize   int64  `json:"fileSize"`
	Ttl        int64  `json:"ttl"`
}

type CreateZapisResponse struct {
	AccessKey string `json:"accessKey"`
	UploadUrl string `json:"uploadUrl"`
}

type GetSaltRequest struct {
	AccessKey string `json:"accessKey"`
}

type GetSaltResponse struct {
	Salt        string `json:"salt"`
	TollPassJwt string `json:"tollpass"`
}

type GetZapisResponse struct {
	DownloadUrl string `json:"downloadUrl"`
	WrappedKey  string `json:"wrappedKey"`
}

type TajnyZapisDynamoItem struct {
	Salt      string `dynamodbav:"salt"`
	AccessKey string `dynamodbav:"access_key"`

	AuthToken  string `dynamodbav:"auth_token"`
	WrappedKey string `dynamodbav:"wrapped_key"`

	S3Key string `dynamodbav:"s3_key"`
	Ttl   int64  `dynamodbav:"ttl"`
}
