package wire

type CreateZapisRequest struct {
	SaltId     string `json:"saltId"`
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
	HashedAccessKey string `json:"accessKey"`
}

type GetSaltResponse struct {
	Salts       map[string]string `json:"salts"`
	TollPassJwt string            `json:"tollpass"`
}

type GetZapisResponse struct {
	DownloadUrl string `json:"downloadUrl"`
	WrappedKey  string `json:"wrappedKey"`
}
