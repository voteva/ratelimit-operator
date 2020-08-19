package ratelimiter

import (
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"strconv"
)

func buildNameForRedis(instance *v1.RateLimiter) string {
	return instance.Name + "-redis"
}

func buildRedisUrl(instance *v1.RateLimiter) string {
	return buildNameForRedis(instance) + ":" + strconv.Itoa(int(constants.REDIS_PORT))
}
