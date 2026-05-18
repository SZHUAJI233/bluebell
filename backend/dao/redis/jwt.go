package redis

import (
	"context"
	"time"
)

func SetToken(token string, tokenExpireDuration time.Duration) error {
	return Set(context.Background(), getRedisKey(KeyJwtTokenSet), token, tokenExpireDuration)
}

func GetToken() (token string, err error) {
	return Get(context.Background(), getRedisKey(KeyJwtTokenSet))
}
