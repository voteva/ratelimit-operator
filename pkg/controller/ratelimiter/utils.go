package ratelimiter

import (
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"strconv"
)

func (r *ReconcileRateLimiter) buildNameForRedis(instance *v1.RateLimiter) string {
	return instance.Name + "-redis"
}

func (r *ReconcileRateLimiter) buildRedisUrl(instance *v1.RateLimiter) string {
	return r.buildNameForRedis(instance) + ":" + strconv.Itoa(int(constants.REDIS_PORT))
}
