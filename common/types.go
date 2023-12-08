package common

type CreateZapisRequest struct {
	Salt       string `json:"salt"`       // uuid4 length 36
	AuthToken  string `json:"authToken"`  // length 44
	WrappedKey string `json:"wrappedKey"` // length 64
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

type VerifyRequest struct {
	AuthToken   string `json:"authToken"`
	TollPassJwt string `json:"tollpass"`
}

type TollPass struct {
	Valid      bool   `json:"valid"`
	Salt       string `json:"salt"`
	WrappedKey string `json:"wrappedKey"`
	AuthToken  string `json:"authToken"`
	S3Key      string `json:"s3Key"`
}

type TajnyZapisDynamoItem struct {
	Salt      string `json:"salt"`
	AccessKey string `json:"accessKey"`

	AuthToken  string `json:"authToken"`
	WrappedKey string `json:"wrappedKey"`

	S3Key string `json:"s3Key"`
	Ttl   int64  `json:"ttl"`
}
