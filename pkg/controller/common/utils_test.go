package common

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
	"github.com/voteva/ratelimit-operator/pkg/utils"
)

func buildRateLimiterConfig(rl *v1.RateLimiter) *v1.RateLimiterConfig {
	host := utils.BuildRandomString(3)
	failureModeDeny := true
	rateLimitRequestTimeout := "0.25s"
	return &v1.RateLimiterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: rl.Namespace,
		},
		Spec: v1.RateLimiterConfigSpec{
			ApplyTo:     v1.GATEWAY,
			Host:        &host,
			Port:        int32(utils.BuildRandomInt(2)),
			RateLimiter: rl.Name,
			Descriptors: []v1.Descriptor{{
				Key: utils.BuildRandomString(3),
			}},
			RateLimitRequestTimeout: &rateLimitRequestTimeout,
			FailureModeDeny:         &failureModeDeny,
		},
	}
}

func buildRateLimiter() *v1.RateLimiter {
	logLevel := v1.INFO
	size := int32(1)

	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
		Spec: v1.RateLimiterSpec{
			LogLevel: &logLevel,
			Size:     &size,
		},
		Status: v1.RateLimiterStatus{},
	}
}
