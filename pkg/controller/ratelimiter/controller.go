package ratelimiter

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

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

var controllerName = "controller_ratelimiter"
var log = logf.Log.WithName(controllerName)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRateLimiter{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.RateLimiter{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRateLimiter{}

type ReconcileRateLimiter struct {
	Client client.Client
	Scheme *runtime.Scheme
}

func (r *ReconcileRateLimiter) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.TODO()

	instance := &v1.RateLimiter{}
	err := r.Client.Get(ctx, request.NamespacedName, instance)
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
