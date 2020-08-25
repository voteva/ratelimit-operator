package ratelimiter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildRedisContainer(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Redis container", func(t *testing.T) {
		containerName := utils.BuildRandomString(3)

		expectedCommand := []string{"redis-server"}
		expectedArgs := []string{"--save", "\"\"", "--appendonly", "no", "--protected-mode", "no", "--bind", "0.0.0.0"}

		actualResult := buildRedisContainer(containerName)

		a.Equal(containerName, actualResult.Name)
		a.Equal(constants.REDIS_IMAGE, actualResult.Image)
		a.Equal(1, len(actualResult.Ports))
		a.Equal(constants.REDIS_PORT, actualResult.Ports[0].ContainerPort)
		a.Equal(corev1.ProtocolTCP, actualResult.Ports[0].Protocol)
		a.Equal(expectedCommand, actualResult.Command)
		a.Equal(expectedArgs, actualResult.Args)
	})
}

func Test_BuildServiceContainer(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ratelimit-service container", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedCommand := []string{"sh", "-c", "/bin/startup.sh"}
		expectedConfigMountPath := fmt.Sprintf("%s/%s/%s", constants.RUNTIME_ROOT, constants.RUNTIME_SUBDIRECTORY, "config")

		actualResult := buildServiceContainer(rateLimiter)

		a.Equal(rateLimiter.Name, actualResult.Name)
		a.Equal(expectedCommand, actualResult.Command)
		a.Equal(constants.RATE_LIMITER_SERVICE_IMAGE, actualResult.Image)
		a.Equal(1, len(actualResult.Ports))
		a.Equal(*rateLimiter.Spec.Port, actualResult.Ports[0].ContainerPort)
		a.Equal(corev1.ProtocolTCP, actualResult.Ports[0].Protocol)
		a.Equal(1, len(actualResult.VolumeMounts))
		a.Equal("config", actualResult.VolumeMounts[0].Name)
		a.Equal(expectedConfigMountPath, actualResult.VolumeMounts[0].MountPath)
	})
}
