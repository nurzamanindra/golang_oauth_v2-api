package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/nurzamanindra/golang_oauth_v2-api/utils/crypto_utils"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(strings.ToLower(at.AccessToken))

	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token")
	}

	return nil
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) GenerateNewExpired() {
	at.Expires = time.Now().UTC().Add(expirationTime * time.Hour).Unix()

}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
