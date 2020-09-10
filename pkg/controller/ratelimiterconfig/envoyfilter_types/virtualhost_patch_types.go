package envoyfilter_types

import v1 "ratelimit-operator/pkg/apis/operators/v1"

type VirtualHostPatchValues struct {
	RateLimits []v1.RateLimits `json:"rate_limits" yaml:"rate_limits"`
}
