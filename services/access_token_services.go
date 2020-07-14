package services

import (
	"github.com/nurzamanindra/golang_oauth_v2-api/datasources/rest"
	"github.com/nurzamanindra/golang_oauth_v2-api/domain/access_token"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
)

var (
	AccessTokenService accessTokenServiceInterface = &accessTokenService{}
)

type accessTokenService struct{}

type accessTokenServiceInterface interface {
	GetTokenById(string) (*access_token.AccessToken, *errors.RestErr)
	CreateAccessToken(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
}

func (at *accessTokenService) GetTokenById(id string) (*access_token.AccessToken, *errors.RestErr) {
	var at_db access_token.AccessToken
	if id == "" {
		return nil, errors.NewBadRequestError("invalid access token id")
	}
	if err := at_db.GetTokenById(id); err != nil {
		return nil, err
	}
	if at_db.IsExpired() {
		return nil, errors.NewNotFoundError("access token for given user id is expired")
	}
	return &at_db, nil
}

func (at *accessTokenService) CreateAccessToken(payload access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	user, err := rest.UserRepository.LoginUser(payload.Username, payload.Password)
	if err != nil {
		return nil, err
	}

	token := access_token.AccessToken{}
	token.GenerateNewExpired()
	token.Generate(user)
	if err := token.SaveToken(); err != nil {
		return nil, err
	}
	return &token, nil
}
