package ratelimiterconfig

import (
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
)

func (r *ReconcileRateLimiterConfig) buildRateLimiterServicePort() int32 {
	return utils.DefaultIfAbsent(r.rateLimiter.Spec.Port, constants.RATELIMITER_PORT)
}
