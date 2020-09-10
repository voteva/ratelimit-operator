package configmapwatcher

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/controller/common"

	v1 "ratelimit-operator/pkg/apis/operators/v1"

	corev1 "k8s.io/api/core/v1"
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

var controllerName = "controller_configmapwatcher"
var log = logf.Log.WithName(controllerName)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileConfigMapWatcher{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(
		&source.Kind{Type: &corev1.ConfigMap{}},
		&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &v1.RateLimiter{}},
		common.CreateOrUpdateConfigMapPredicate)

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileConfigMapWatcher{}

type ReconcileConfigMapWatcher struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileConfigMapWatcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
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

	return reconcile.Result{}, nil
}

func (r *ReconcileConfigMapWatcher) reconcileConfigMap(ctx context.Context, instance *v1.RateLimiter) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Name", instance.Name)

	foundConfigMap := &corev1.ConfigMap{}

	err := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		} else {
			reqLogger.Error(err, "Failed to get ConfigMap")
			return reconcile.Result{}, err
		}
	}

	list, err := r.getRateLimiterLists(ctx, instance)
	if err != nil || list == nil {
		return reconcile.Result{}, err
	}

	r.updateConfigMap(ctx, foundConfigMap, list)

	return reconcile.Result{}, nil
}

func (r *ReconcileConfigMapWatcher) getRateLimiterLists(ctx context.Context, instance *v1.RateLimiter) (*v1.RateLimiterConfigList, error) {
	opts := []client.ListOption{client.InNamespace(instance.Namespace)}
	list := &v1.RateLimiterConfigList{}
	err := r.client.List(ctx, list, opts...)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return list, err
}

func (r *ReconcileConfigMapWatcher) updateConfigMap(ctx context.Context, configMap *corev1.ConfigMap, list *v1.RateLimiterConfigList) {
	data := configMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	needUpdate := false
	for _, rlc := range list.Items {
		fileName := common.BuildConfigMapDataFileName(rlc.Name)
		expectedVal := common.BuildRateLimitPropertyValue(&rlc)
		actualVal, found := data[fileName]

		if !found || actualVal != expectedVal {
			data[fileName] = expectedVal
			needUpdate = true
		}
	}

	if needUpdate {
		configMap.Data = data
		r.client.Update(ctx, configMap)
	}
}
