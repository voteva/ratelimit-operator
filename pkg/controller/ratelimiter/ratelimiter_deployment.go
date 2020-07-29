package ratelimiter

import (
	"context"
	"ratelimit-operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeploymentForService(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundDeployment := &appsv1.Deployment{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.buildDeploymentForService(instance)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildDeploymentForService(instance *v1.RateLimiter) *appsv1.Deployment {
	labels := utils.LabelsForApp(instance.Name)
	defaultRedisUrl := r.buildRedisUrl(instance)
	replicas := int32(2)
	var defaultMode int32 = 420

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: utils.AnnotationSidecarIstio(),
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: instance.Name,
								},
								DefaultMode: &defaultMode,
							},
						},
					}},
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: "redis:alpine",
						},
						{
							Name:  instance.Name,
							Image: instance.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									ContainerPort: r.buildRateLimiterServicePort(instance),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "LOG_LEVEL",
									Value: utils.DefaultIfEmpty(string(*instance.Spec.LogLevel), string(v1.INFO)),
								},
								{
									Name:  "REDIS_SOCKET_TYPE",
									Value: "TCP",
								},
								{
									Name:  "REDIS_URL",
									Value: defaultRedisUrl,
								},
								{
									Name:  "RUNTIME_IGNOREDOTFILES",
									Value: "true",
								},
								{
									Name:  "RUNTIME_ROOT",
									Value: "/home/user/src/runtime/data",
								},
								{
									Name:  "RUNTIME_SUBDIRECTORY",
									Value: "ratelimit",
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
								MountPath: "/home/user/src/runtime/data/ratelimit/config",
							}},
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							EnvFrom: []corev1.EnvFromSource{{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.Name,
									},
								},
							}},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, dep, r.scheme)
	return dep
}
