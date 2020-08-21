package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildServiceForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Service for Redis", func(t *testing.T) {
		rateLimiter := &v1.RateLimiter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.BuildRandomString(3),
				Namespace: utils.BuildRandomString(3),
			},
		}

		expectedServiceName := buildNameForRedis(rateLimiter.Name)

		actualResult := buildServiceForRedis(rateLimiter)

		a.Equal(expectedServiceName, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(1, len(actualResult.Spec.Ports))
		a.Equal(expectedServiceName, actualResult.Spec.Ports[0].Name)
		a.Equal(corev1.ProtocolTCP, actualResult.Spec.Ports[0].Protocol)
		a.Equal(constants.REDIS_PORT, actualResult.Spec.Ports[0].Port)
		a.Equal(intstr.IntOrString{Type: intstr.Int, IntVal: constants.REDIS_PORT}, actualResult.Spec.Ports[0].TargetPort)
		a.Equal(utils.SelectorsForApp(expectedServiceName), actualResult.Spec.Selector)
	})
}
