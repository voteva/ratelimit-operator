package common

import (
	"github.com/ghodss/yaml"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/common/types"
)

func BuildRateLimitPropertyValue(instance *v1.RateLimiterConfig) string {
	prop := &types.RateLimitProperty{
		Domain:      instance.Name,
		Descriptors: instance.Spec.Descriptors,
	}

	res, _ := yaml.Marshal(prop)
	return string(res)
}

func BuildConfigMapDataFileName(name string) string {
	return name + ".yaml"
}
