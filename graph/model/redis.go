package model

import (
	"errors"
	"time"

	"gitlab.com/sirinibin/go-mysql-graphql/config"
)

func (token *Token) SaveToRedis() error {
	expires := time.Unix(int64(token.ExpiresAt), 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	errAccess := config.RedisClient.Set(token.AccessUUID, token.UserID, expires.Sub(now)).Err()

	return errAccess
}

func (token *Token) ExistsInRedis() error {

	userid, err := config.RedisClient.Get(token.AccessUUID).Result()
	if err != nil {
		return err
	}

	if token.UserID != userid {
		return errors.New("User id doesn't exist in redis!")
	}

	return nil

}
