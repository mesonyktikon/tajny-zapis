package tokens

import (
	"encoding/hex"
	"os"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type TollPass struct {
	Valid      bool   `json:"valid"`
	WrappedKey string `json:"wrappedKey"`
	AuthToken  string `json:"authToken"`
	S3Key      string `json:"s3Key"`
}

var key []byte
var enc jose.Encrypter

func init() {
	secret := os.Getenv("TOLLPASS_SECRET")
	if len(secret) != 64 {
		panic("TOLLPASS_SECRET must be 64 characters long")
	}

	key = make([]byte, 32)
	_, err := hex.Decode(key, []byte(secret))
	if err != nil {
		panic(err)
	}

	jwk := jose.JSONWebKey{
		Key:       key,
		KeyID:     "v1",
		Algorithm: "dir",
		Use:       "enc",
	}

	enc_, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: jwk},
		(&jose.EncrypterOptions{}).WithType("JWT"),
	)
	if err != nil {
		panic(err)
	}
	enc = enc_
}

func GenerateTollPassJwt(tollPass *TollPass) (string, error) {
	now := jwt.NumericDate(time.Now().Unix())
	exp := jwt.NumericDate(time.Now().Add(time.Minute).Unix())
	claims := jwt.Claims{
		Subject:   "tollpass",
		Issuer:    "tajny-zapis",
		NotBefore: &now,
		Expiry:    &exp,
	}
	return jwt.Encrypted(enc).Claims(claims).Claims(tollPass).CompactSerialize()
}

func DecryptTollPassJwt(token string) (*TollPass, error) {
	t, err := jwt.ParseEncrypted(token)
	if err != nil {
		return nil, err
	}

	var tollPass TollPass
	err = t.Claims(key, &tollPass)
	return &tollPass, err
}
