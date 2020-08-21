package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildService(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Service for ratelimit-service", func(t *testing.T) {
		logLevel := v1.INFO
		port := int32(utils.BuildRandomInt(2))
		size := int32(utils.BuildRandomInt(1))

		rateLimiter := &v1.RateLimiter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.BuildRandomString(3),
				Namespace: utils.BuildRandomString(3),
			},
			Spec: v1.RateLimiterSpec{
				LogLevel: &logLevel,
				Port:     &port,
				Size:     &size,
			},
		}

		actualResult := buildService(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(utils.SelectorsForApp(rateLimiter.Name), actualResult.Spec.Selector)
		a.Equal(1, len(actualResult.Spec.Ports))
		a.Equal("grpc-"+rateLimiter.Name, actualResult.Spec.Ports[0].Name)
		a.Equal(corev1.ProtocolTCP, actualResult.Spec.Ports[0].Protocol)
		a.Equal(*rateLimiter.Spec.Port, actualResult.Spec.Ports[0].Port)
		a.Equal(intstr.IntOrString{Type: intstr.Int, IntVal: *rateLimiter.Spec.Port}, actualResult.Spec.Ports[0].TargetPort)
	})
}
