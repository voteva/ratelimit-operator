package ratelimiter

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
)

func (r *ReconcileRateLimiter) BuildRedisContainer(name string) corev1.Container {
	return corev1.Container{
		Name:  name,
		Image: constants.REDIS_IMAGE,
		Ports: []corev1.ContainerPort{{
			ContainerPort: constants.REDIS_PORT,
			Protocol:      corev1.ProtocolTCP,
		}},
		Command: []string{"redis-server"},
		Args: []string{
			"--save", "\"\"",
			"--appendonly", "no",
			"--protected-mode", "no",
			"--bind", "0.0.0.0",
		},
	}
}

func (r *ReconcileRateLimiter) BuildServiceContainer(instance *v1.RateLimiter) corev1.Container {
	redisUrl := r.buildRedisUrl(instance)
	configMountPath := fmt.Sprintf("%s/%s/%s", constants.RUNTIME_ROOT, constants.RUNTIME_SUBDIRECTORY, "config")

	return corev1.Container{
		Name: instance.Name,
		Command: []string{
			"sh",
			"-c",
			"/bin/startup.sh",
		},
		Image: constants.RATE_LIMITER_SERVICE_IMAGE,
		Ports: []corev1.ContainerPort{{
			ContainerPort: *instance.Spec.Port,
			Protocol:      corev1.ProtocolTCP,
		}},
		Env: []corev1.EnvVar{
			{
				Name:  "LOG_LEVEL",
				Value: string(*instance.Spec.LogLevel),
			},
			{
				Name:  "REDIS_SOCKET_TYPE",
				Value: "TCP",
			},
			{
				Name:  "REDIS_URL",
				Value: redisUrl,
			},
			{
				Name:  "RUNTIME_IGNOREDOTFILES",
				Value: "true",
			},
			{
				Name:  "RUNTIME_ROOT",
				Value: constants.RUNTIME_ROOT,
			},
			{
				Name:  "RUNTIME_SUBDIRECTORY",
				Value: constants.RUNTIME_SUBDIRECTORY,
			},
			{
				Name:  "RUNTIME_WATCH_ROOT",
				Value: "false",
			},
			{
				Name:  "USE_STATSD",
				Value: "false",
			},
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "config",
			MountPath: configMountPath,
		}},
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
	}
}
