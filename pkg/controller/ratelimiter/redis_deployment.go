package ratelimiter

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeploymentForRedis(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundDeployment := &appsv1.Deployment{}
	deploymentName := r.buildNameForRedis(instance)
	deploymentFromInstance := r.buildDeploymentForRedis(instance, deploymentName)

	err := r.Client.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: instance.Namespace}, foundDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment Redis")
			err = r.Client.Create(ctx, deploymentFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment Redis")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Deployment Redis")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundDeployment.Spec, deploymentFromInstance.Spec) {
		foundDeployment.Spec = deploymentFromInstance.Spec
		r.Client.Update(ctx, foundDeployment)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildDeploymentForRedis(instance *v1.RateLimiter, deploymentName string) *appsv1.Deployment {
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
					Labels:      labels,
					Annotations: utils.AnnotationSidecarIstio(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						r.BuildRedisContainer(deploymentName),
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}
