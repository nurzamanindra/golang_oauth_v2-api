package access_token

import (
	access_token_db "github.com/nurzamanindra/golang_oauth_v2-api/datasources/mysql"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/errors"
	"github.com/nurzamanindra/golang_oauth_v2-api/utils/mysql_utils"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_token WHERE access_token=?;"
	queryInsertAccessToken = "INSERT INTO access_token(access_token, user_id, client_id, expires) VALUES(?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_token SET expires=? WHERE access_token=?;"
)

func (at *AccessToken) GetTokenById(access_token_id string) *errors.RestErr {
	stmt, err := access_token_db.Client.Prepare(queryGetAccessToken)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(access_token_id)

	if geterr := result.Scan(&at.AccessToken, &at.UserId, &at.ClientId, &at.Expires); err != nil {
		return mysql_utils.ParseError(geterr)
	}
	return nil

}

func (at *AccessToken) SaveToken() *errors.RestErr {
	stmt, err := access_token_db.Client.Prepare(queryInsertAccessToken)
	if err != nil {
		return errors.NewInternalServerError("error when preparing insert to database")
	}
	defer stmt.Close()

	_, inserterr := stmt.Exec(at.AccessToken, at.UserId, at.ClientId, at.Expires)
	if inserterr != nil {
		return errors.NewInternalServerError("error when trying to insert to database")
	}

	return nil
}

func (at *AccessToken) UpdateTokenExpired() *errors.RestErr {
	return nil
}
