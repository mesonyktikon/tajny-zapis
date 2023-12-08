package logic

import (
	"encoding/hex"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"tuffbizz.com/m/v2/common"
)

const secret = "7370c6bce277bd1218a61cd30adc47c96174a0c55628bc73e9c3f94202e2e377"

var key []byte

var enc jose.Encrypter

func init() {
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

func GenerateTollPassJwt(tollPass *common.TollPass) (string, error) {
	now := jwt.NumericDate(time.Now().Unix())
	exp := jwt.NumericDate(time.Now().Add(time.Minute).Unix())
	claims := jwt.Claims{
		Subject:   "tollpass",
		Issuer:    "tajny-zapis",
		NotBefore: &now,
		Expiry:    &exp,
	}

	token, err := jwt.Encrypted(enc).Claims(claims).Claims(tollPass).CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}

func DecryptTollPassJwt(token string) (*common.TollPass, error) {
	t, err := jwt.ParseEncrypted(token)
	if err != nil {
		return nil, err
	}

	var tollPass common.TollPass
	err = t.Claims(key, &tollPass)
	if err != nil {
		return nil, err
	}

	return &tollPass, nil
}
