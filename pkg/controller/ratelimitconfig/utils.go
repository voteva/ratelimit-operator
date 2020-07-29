package ratelimitconfig

import (
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
)

func (r *ReconcileRateLimitConfig) buildRateLimiterServiceFqdn() string {
	return r.rateLimiter.Name + "." + r.rateLimiter.Namespace + ".svc.cluster.local"
}

func (r *ReconcileRateLimitConfig) buildRateLimiterServicePort() int32 {
	return utils.DefaultIfAbsent(r.rateLimiter.Spec.Port, constants.DEFAULT_RATELIMITER_PORT)
}
