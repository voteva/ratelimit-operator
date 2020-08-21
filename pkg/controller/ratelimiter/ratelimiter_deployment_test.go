package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildDeployment(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Deployment for ratelimit-service", func(t *testing.T) {
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

		actualResult := buildDeploymentForService(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(rateLimiter.Spec.Size, actualResult.Spec.Replicas)
		a.Equal(utils.LabelsForApp(rateLimiter.Name), actualResult.Spec.Selector.MatchLabels)
		a.Equal(utils.LabelsForApp(rateLimiter.Name), actualResult.Spec.Template.ObjectMeta.Labels)
		a.Equal(utils.AnnotationSidecarIstio(), actualResult.Spec.Template.ObjectMeta.Annotations)
		a.Equal(1, len(actualResult.Spec.Template.Spec.Containers))
		a.Equal(1, len(actualResult.Spec.Template.Spec.Volumes))
		a.Equal("config", actualResult.Spec.Template.Spec.Volumes[0].Name)
		a.Equal(rateLimiter.Name, actualResult.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name)
	})
}
