package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildConfigMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ConfigMap for ratelimit-service", func(t *testing.T) {
		rateLimiter := &v1.RateLimiter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.BuildRandomString(3),
				Namespace: utils.BuildRandomString(3),
			},
		}

		actualResult := buildConfigMap(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(map[string]string{}, actualResult.Data)
	})
}
