package ratelimiterconfig

import (
	"context"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var controllerName = "controller_ratelimiter_config"
var log = logf.Log.WithName(controllerName)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRateLimiterConfig{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.RateLimiterConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1alpha3.EnvoyFilter{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1.RateLimiterConfig{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRateLimiterConfig{}

type ReconcileRateLimiterConfig struct {
	client      client.Client
	scheme      *runtime.Scheme
	rateLimiter *v1.RateLimiter
	configMap   *corev1.ConfigMap
}

func (r *ReconcileRateLimiterConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)
	ctx := context.TODO()

	instance := &v1.RateLimiterConfig{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if isNeedUpdateWithDefaults(instance) {
		r.client.Update(ctx, instance)
	}

	err = r.getRateLimiter(ctx, instance)
	if err != nil {
		reqLogger.Error(err, "Error get RateLimiter [%s]", instance.Spec.RateLimiter)
		return reconcile.Result{}, err
	}

	err = r.getRateLimiterConfigMap(ctx, instance)
	if err != nil {
		reqLogger.Error(err, "Error get RateLimiter ConfigMap [%s]", instance.Spec.RateLimiter)
		return reconcile.Result{}, err
	}

	if result, err := r.updateConfigMap(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	if result, err := r.reconcileEnvoyFilter(ctx, instance); err != nil || result.Requeue {
		return result, err
	}

	r.addFinalizerIfNotExists(ctx, instance)

	if utils.IsBeingDeleted(instance) {
		if !utils.HasFinalizer(instance, controllerName) {
			return reconcile.Result{}, nil
		}
		if err := r.manageCleanUpLogic(ctx, instance); err != nil {
			return reconcile.Result{}, err
		}
		utils.RemoveFinalizer(instance, controllerName)
		if err := r.client.Update(ctx, instance); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiterConfig) getRateLimiter(ctx context.Context, instance *v1.RateLimiterConfig) error {
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

func (r *ReconcileRateLimiterConfig) getRateLimiterConfigMap(ctx context.Context, instance *v1.RateLimiterConfig) error {
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

func (r *ReconcileRateLimiterConfig) addFinalizerIfNotExists(ctx context.Context, instance *v1.RateLimiterConfig) {
	if !utils.HasFinalizer(instance, controllerName) {
		utils.AddFinalizer(instance, controllerName)
		r.client.Update(ctx, instance)
	}
}

func (r *ReconcileRateLimiterConfig) manageCleanUpLogic(context context.Context, instance *v1.RateLimiterConfig) error {
	if err := r.deleteFromConfigMap(context, instance); err != nil {
		log.Error(err, "Failed to clean up ConfigMap for config [%s]", instance.Name)
		return err
	}
	return nil
}

func isNeedUpdateWithDefaults(instance *v1.RateLimiterConfig) bool {
	needUpdate := false

	if instance.Spec.FailureModeDeny == nil {
		defaultFailureModeDeny := false
		instance.Spec.FailureModeDeny = &defaultFailureModeDeny
		needUpdate = true
	}
	if instance.Spec.RateLimitRequestTimeout == nil {
		defaultRateLimitRequestTimeout := "0.25s"
		instance.Spec.RateLimitRequestTimeout = &defaultRateLimitRequestTimeout
		needUpdate = true
	}
	return needUpdate
}
