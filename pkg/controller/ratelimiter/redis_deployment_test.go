package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildDeploymentForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Deployment for Redis", func(t *testing.T) {
		rateLimiter := &v1.RateLimiter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.BuildRandomString(3),
				Namespace: utils.BuildRandomString(3),
			},
		}

		expectedDeploymentName := buildNameForRedis(rateLimiter.Name)
		expectedLabels := utils.LabelsForApp(expectedDeploymentName)
		expectedReplicas := int32(1)

		actualResult := buildDeploymentForRedis(rateLimiter)

		a.Equal(expectedDeploymentName, actualResult.Name)
		a.Equal(rateLimiter.Namespace, actualResult.Namespace)
		a.Equal(&expectedReplicas, actualResult.Spec.Replicas)
		a.Equal(expectedLabels, actualResult.Spec.Selector.MatchLabels)
		a.Equal(expectedLabels, actualResult.Spec.Template.ObjectMeta.Labels)
		a.Equal(utils.AnnotationSidecarIstio(), actualResult.Spec.Template.ObjectMeta.Annotations)
		a.Equal(1, len(actualResult.Spec.Template.Spec.Containers))
	})
}
