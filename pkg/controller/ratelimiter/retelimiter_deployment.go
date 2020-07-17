package ratelimiter

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeployment(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundDeployment := &appsv1.Deployment{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.buildDeployment(instance)
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

	size := instance.Spec.Size
	if *foundDeployment.Spec.Replicas != size {
		foundDeployment.Spec.Replicas = &size
		err = r.client.Update(ctx, foundDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildDeployment(instance *v1.RateLimiter) *appsv1.Deployment {
	ls := LabelsForRateLimiter(instance.Name)
	var defaultMode int32 = 420 // TODO

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Size,
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
									Name: instance.Name,
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
							Name:  instance.Name,
							Image: "evil26r/service_rite_limit",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									ContainerPort: instance.Spec.ServicePort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "config",
								MountPath: "/data/ratelimit/config",
							}},
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							EnvFrom: []corev1.EnvFromSource{{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.Name,
									},
								},
							}},
						}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, dep, r.scheme)
	return dep
}
