package ratelimitconfig

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	err = c.Watch(&source.Kind{Type: &v1.RateLimitConfig{}}, &handler.EnqueueRequestForObject{})
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
	configMap   *corev1.ConfigMap
}

func (r *ReconcileRateLimitConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling RateLimitConfig")

	ctx := context.TODO()

	instance := &v1.RateLimitConfig{}
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

	err = r.getRateLimiterConfigMap(ctx, instance)
	if err != nil {
		reqLogger.Error(err, "Get RateLimiter ConfigMap [%s/%s] error", instance.Spec.RateLimiter, instance.Namespace)
		return reconcile.Result{}, err
	}

	if result, err := r.reconcileConfigMap(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileEnvoyFilter(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimitConfig) getRateLimiter(ctx context.Context, instance *v1.RateLimitConfig) error {
	rateLimiter := &v1.RateLimiter{}
	err := r.client.Get(
		ctx,
		types.NamespacedName{
			Name:      instance.Spec.RateLimiter,
			Namespace: instance.Namespace,
		},
		rateLimiter,
	)
	if err != nil {
		return err
	}
	r.rateLimiter = rateLimiter
	return nil
}

func (r *ReconcileRateLimitConfig) getRateLimiterConfigMap(ctx context.Context, instance *v1.RateLimitConfig) error {
	configMap := &corev1.ConfigMap{}
	err := r.client.Get(
		ctx,
		types.NamespacedName{
			Name:      instance.Spec.RateLimiter,
			Namespace: instance.Namespace,
		},
		configMap,
	)
	if err != nil {
		return err
	}
	r.configMap = configMap
	return nil
}
