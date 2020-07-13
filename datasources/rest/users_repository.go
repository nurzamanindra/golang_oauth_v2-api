package rest

import (
	"encoding/json"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/nurzamanindra/golang_oauth_v2-api/domain/users"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
)

var (
	UserRepository RestUserRepository = &userRepository{}
	client                            = resty.New()
)

type RestUserRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type userRepository struct{}

func (u *userRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	client.SetHostURL("https://localhost:9000").SetTimeout(1000 * time.Millisecond)
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := client.R().SetBody(request).Post("/users/login")

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	if response.Size() <= 2 {
		return nil, errors.NewInternalServerError("invalid rest client response when trying get user")
	}
	if response.StatusCode() > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(response.Body(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}
	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, errors.NewNotFoundError("error when trying unmarshal users response")
	}
	return &user, nil
}
