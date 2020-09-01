package ratelimiter

import (
	"context"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeploymentForService(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundDeployment := &appsv1.Deployment{}
	deploymentFromInstance := buildDeploymentForService(instance)
	_ = controllerutil.SetControllerReference(instance, deploymentFromInstance, r.scheme)

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment")
			err = r.client.Create(ctx, deploymentFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get Deployment")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundDeployment.Spec, deploymentFromInstance.Spec) {
		foundDeployment.Spec = deploymentFromInstance.Spec
		r.client.Update(ctx, foundDeployment)
	}

	return reconcile.Result{}, nil
}

func buildDeploymentForService(instance *v1.RateLimiter) *appsv1.Deployment {
	labels := utils.LabelsForApp(instance.Name)
	var defaultMode int32 = 420

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: utils.Merge(utils.AnnotationSidecarIstio(), utils.AnnotationMetricsIstio(9102)),
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
					},
						{
							Name: "config-statsd-exporter",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "rate-limiter-statsd-exporter",
									},
									DefaultMode: &defaultMode,
								}},
						},
					},
					Containers: []corev1.Container{
						buildServiceContainer(instance),
						buildStatsdExporterContainer("statsd-exporter"),
					},
				},
			},
		},
	}
	return dep
}
