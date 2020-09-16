package types

import v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"

type RateLimitProperty struct {
	Domain      string          `json:"domain" yaml:"domain"`
	Descriptors []v1.Descriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}
