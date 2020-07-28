package ratelimiter

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"

	"ratelimit-operator/pkg/apis/operators/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ratelimiter")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRateLimiter{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("ratelimiter-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource RateLimiter
	err = c.Watch(&source.Kind{Type: &v1.RateLimiter{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//Watch for changes to secondary resources and requeue the owner RateLimiter
	log.Info("Watch for changes to appsv1.Deployment")
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	//log.Info("Watch for changes to corev1.Service")
	//err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &v1.RateLimiter{},
	//})
	//if err != nil {
	//	return err
	//}

	//log.Info("Watch for changes to corev1.ConfigMap")
	//err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &v1.RateLimiter{},
	//})
	//if err != nil {
	//	return err
	//}

	//log.Info("Watch for changes to v1alpha3.EnvoyFilter")
	//err = c.Watch(&source.Kind{Type: &v1alpha3.EnvoyFilter{}},
	//	&handler.EnqueueRequestForOwner{
	//		IsController: true,
	//		OwnerType:    &v1.RateLimiter{},
	//	})
	//if err != nil {
	//	return err
	//}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRateLimiter{}

type ReconcileRateLimiter struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileRateLimiter) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling RateLimiter")

	ctx := context.TODO()

	instance := &v1.RateLimiter{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if result, err := r.reconcileConfigMap(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileDeploymentForRedis(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileServiceForRedis(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileDeploymentForService(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileServiceForService(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	return reconcile.Result{}, nil
}
