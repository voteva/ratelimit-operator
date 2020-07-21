package ratelimitconfig

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	operatorsv1 "ratelimit-operator/pkg/apis/operators/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ratelimitconfig")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRateLimitConfig{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("ratelimitconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &operatorsv1.RateLimitConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRateLimitConfig{}

type ReconcileRateLimitConfig struct {
	client      client.Client
	scheme      *runtime.Scheme
	rateLimiter *v1.RateLimiter
}

func (r *ReconcileRateLimitConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling RateLimitConfig")

	ctx := context.TODO()

	instance := &operatorsv1.RateLimitConfig{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	err = r.getRateLimiter(ctx, instance)
	if err != nil {
		reqLogger.Error(err, "Get RateLimiter[%s/%s] error", instance.Spec.RateLimiter, instance.Namespace)
		return reconcile.Result{}, err
	}

	reqLogger.Info("Success")

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimitConfig) getRateLimiter(ctx context.Context, instance *operatorsv1.RateLimitConfig) error {
	rateLimiter := &v1.RateLimiter{}
	err := r.client.Get(
		ctx,
		types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      instance.Spec.RateLimiter,
		},
		rateLimiter,
	)
	if err != nil {
		return err
	}
	r.rateLimiter = rateLimiter
	return nil
}
