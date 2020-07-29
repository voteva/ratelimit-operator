package ratelimiter

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeploymentForRedis(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	foundDeployment := &appsv1.Deployment{}
	deploymentName := r.buildNameForRedis(instance)

	err := r.client.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.buildDeploymentForRedis(instance)
		log.Info("Creating a new Deployment Redis", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment Redis", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment Redis")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildDeploymentForRedis(instance *v1.RateLimiter) *appsv1.Deployment {
	deploymentName := r.buildNameForRedis(instance)
	labels := utils.LabelsForApp(deploymentName)
	replicas := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: r.buildRedisImage(instance),
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, dep, r.scheme)
	return dep
}
