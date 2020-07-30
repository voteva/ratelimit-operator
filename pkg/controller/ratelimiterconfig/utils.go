package ratelimiterconfig

import (
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
)

func (r *ReconcileRateLimiterConfig) buildRateLimiterServiceFqdn() string {
	return r.rateLimiter.Name + "." + r.rateLimiter.Namespace + ".svc.cluster.local"
}

func (r *ReconcileRateLimiterConfig) buildRateLimiterServicePort() int32 {
	return utils.DefaultIfAbsent(r.rateLimiter.Spec.Port, constants.DEFAULT_RATELIMITER_PORT)
}
