package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
	"github.com/stretchr/testify/assert"
)

type mockRestErr struct {
	Message_error string `json:"message_error"`
	Status        string `json:"status"`
	Error         string `json:"error"`
}

type mockUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases")
	os.Exit(m.Run())
}

func TestLoginUserTimeOutAPI(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("POST", "https://localhost:9000/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(500, ``)
			if err != nil {
				return httpmock.NewStringResponse(500, `{}`), nil
			}
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	repo := userRepository{}

	user, err := repo.LoginUser("test@gmail.com", "hahaha123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client response when trying get user", err.Message)
}

func TestInvalidErrorInterface(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())

	mockError := mockRestErr{
		Message_error: "invalid login credential",
		Status:        "500",
		Error:         "internal_server_error",
	}

	httpmock.RegisterResponder("POST", "https://localhost:9000/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(500, mockError)
			if err != nil {
				return httpmock.NewStringResponse(500, `{}`), nil
			}
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	repo := userRepository{}

	user, err := repo.LoginUser("test@gmail.com", "hahaha123")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)

}

func TestInvalidLoginCredential(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())

	mockError := errors.RestErr{
		Message: "invalid login credential",
		Status:  404,
		Error:   "not_found",
	}
	httpmock.RegisterResponder("POST", "https://localhost:9000/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(404, mockError)
			if err != nil {
				return httpmock.NewStringResponse(404, ``), nil
			}
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	repo := userRepository{}

	user, err := repo.LoginUser("test@gmail.com", "hahaha123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credential", err.Message)

}

func TestLoginUserJsonResponse(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("POST", "https://localhost:9000/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, `mcUser`)
			if err != nil {
				return httpmock.NewStringResponse(200, ``), nil
			}
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	repo := userRepository{}

	user, err := repo.LoginUser("test@gmail.com", "hahaha123")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "error when trying unmarshal users response", err.Message)

}

func TestLoginUserNoError(t *testing.T) {
	httpmock.Activate()
	httpmock.ActivateNonDefault(client.GetClient())
	mcUser := mockUser{
		Id:        123,
		FirstName: "mock",
		LastName:  "mock",
		Email:     "mock@mock.com",
	}
	httpmock.RegisterResponder("POST", "https://localhost:9000/users/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, mcUser)
			if err != nil {
				return httpmock.NewStringResponse(200, ``), nil
			}
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	repo := userRepository{}

	user, err := repo.LoginUser("test@gmail.com", "hahaha123")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 123, user.Id)
	assert.EqualValues(t, "mock", user.FirstName)
	assert.EqualValues(t, "mock", user.LastName)
	assert.EqualValues(t, "mock@mock.com", user.Email)
}