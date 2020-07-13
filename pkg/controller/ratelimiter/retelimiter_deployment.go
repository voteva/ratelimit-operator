package ratelimiter

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileRateLimiter) DeploymentForRateLimiter(m *operatorsv1alpha1.RateLimiter) *appsv1.Deployment {
	ls := LabelsForRateLimiter(m.Name)
	var replicas int32 = 1      // TODO
	var defaultMode int32 = 420 // TODO

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ls,
					Annotations: map[string]string{"sidecar.istio.io/inject": "true"},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: m.Name,
								},
								DefaultMode: &defaultMode,
								Items: []corev1.KeyToPath{{
									Key:  "rate_limit.property",
									Path: "config.yaml",
								}},
							},
						},
					}},
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: "redis:alpine",
						},
						{
							Name:  "rate-limit-server",
							Image: "evil26r/service_rite_limit",
							Ports: []corev1.ContainerPort{{
								ContainerPort: 8080,
								Protocol:      "TCP",
							}},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "config",
								MountPath: "/data/ratelimit/config",
							}},
							TerminationMessagePolicy: "File",
							EnvFrom: []corev1.EnvFromSource{{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: m.Name,
									},
								},
							}},
						}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}
