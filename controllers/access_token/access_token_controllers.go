package access_token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurzamanindra/golang_oauth_v2-api/domain/access_token"
	"github.com/nurzamanindra/golang_oauth_v2-api/services"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
)

func GetAccessTokenById(c *gin.Context) {
	at_id := c.Param("access_token_id")
	result, err := services.AccessTokenService.GetTokenById(at_id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func CreateAccessToken(c *gin.Context) {
	var payload access_token.AccessToken
	if err := c.ShouldBindJSON(&payload); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	err := services.AccessTokenService.CreateAccessToken(&payload)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, payload)
}
