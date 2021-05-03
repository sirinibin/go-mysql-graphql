package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jameskeane/bcrypt"
)

func AuthenticateByAuthCode(tokenStr string) (tokenClaims TokenClaims, err error) {

	if govalidator.IsNull(tokenStr) {
		return tokenClaims, errors.New("auth_code is required")
	}

	tokenClaims, err = AuthenticateByJWTToken(tokenStr)
	if err != nil {
		return tokenClaims, err
	}
	if tokenClaims.Type != "auth_code" {
		return tokenClaims, errors.New("Invalid auth code")
	}

	return tokenClaims, nil

}

//GenerateAuthCode : generate and return authcode
func (auth *AuthorizeRequest) GenerateAuthCode() (authCode AuthCodeResponse, err error) {

	// Generate Auth code
	expiresAt := time.Now().Add(time.Hour * 5) // expiry for auth code is 5min
	token, err := generateAndSaveToken(auth.Username, expiresAt, "auth_code")
	authCode.ExpiresAt = int(token.ExpiresAt)
	authCode.Code = token.TokenStr

	return authCode, err
}

//Validate : Validate authorization data
func (auth *AuthorizeRequest) Authenticate() error {

	if govalidator.IsNull(auth.Username) {
		return errors.New("Username is required")
	}
	if govalidator.IsNull(auth.Password) {
		return errors.New("Password is required")
	}

	user, err := FindUserByUsername(auth.Username)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows || !bcrypt.Match(auth.Password, user.Password) {
		return errors.New("Username or Password is wrong")
	}

	return nil
}
