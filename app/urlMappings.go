package app

import "github.com/nurzamanindra/golang_oauth_v2-api/controllers/access_token"

func mapUrls() {
	router.GET("/oauth/access_token/:access_token_id", access_token.GetAccessTokenById)
	router.POST("/oauth/access_token", access_token.CreateAccessToken)
}
