package ratelimiter

import (
	"ratelimit-operator/pkg/constants"
	"strconv"
)

func buildNameForRedis(name string) string {
	return name + "-redis"
}

func buildRedisUrl(name string) string {
	return buildNameForRedis(name) + ":" + strconv.Itoa(int(constants.REDIS_PORT))
}
