package ratelimiter

import (
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"strconv"
)

func (r *ReconcileRateLimiter) buildNameForRedis(instance *v1.RateLimiter) string {
	return instance.Name + "-redis"
}

func (r *ReconcileRateLimiter) buildRedisUrl() string {
	return "localhost" + ":" + strconv.Itoa(int(constants.DEFAULT_REDIS_PORT))
}

func (r *ReconcileRateLimiter) buildRateLimiterServicePort(instance *v1.RateLimiter) int32 {
	return utils.DefaultIfAbsent(instance.Spec.Port, int32(constants.DEFAULT_RATELIMITER_PORT))
}
