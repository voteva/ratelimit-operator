package ratelimiter

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiter) reconcileDeployment(request reconcile.Request, instance *operatorsv1alpha1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	foundDeployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.buildDeployment(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// TODO доработать
	// Check Pod size = 1
	var expectedSize int32 = 1
	if *foundDeployment.Spec.Replicas != expectedSize {
		foundDeployment.Spec.Replicas = &expectedSize
		err = r.client.Update(context.TODO(), foundDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", foundDeployment.Namespace, "Pod.Name", foundDeployment.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiter) buildDeployment(m *operatorsv1alpha1.RateLimiter) *appsv1.Deployment {
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
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8081,
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
