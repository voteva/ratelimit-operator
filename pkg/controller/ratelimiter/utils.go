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

func (r *ReconcileRateLimiter) buildRedisUrl(instance *v1.RateLimiter) string {
	return r.buildNameForRedis(instance) + ":" + strconv.Itoa(int(constants.REDIS_PORT))
	//return "localhost" + ":" + strconv.Itoa(int(constants.REDIS_PORT))
}

func (r *ReconcileRateLimiter) buildRateLimiterServicePort(instance *v1.RateLimiter) int32 {
	return utils.DefaultIfAbsent(instance.Spec.Port, int32(constants.RATELIMITER_PORT))
}
