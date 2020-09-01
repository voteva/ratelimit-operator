package configmapwatcher

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/controller/common"

	"github.com/ghodss/yaml"
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

	opts := []client.ListOption{client.InNamespace(instance.Namespace)}
	list := &v1.RateLimiterConfigList{}
	err = r.client.List(ctx, list, opts...)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	data := foundConfigMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	needUpdate := false
	for _, rlc := range list.Items {
		fileName := buildConfigMapDataFileName(rlc.Name)
		expectedVal := buildRateLimitPropertyValue(rlc.Spec.RateLimitProperty)
		actualVal, found := data[fileName]

		if !found || actualVal != expectedVal {
			data[fileName] = expectedVal
			needUpdate = true
		}
	}

	if needUpdate {
		foundConfigMap.Data = data
		r.client.Update(ctx, foundConfigMap)
	}

	return reconcile.Result{}, nil
}

func buildRateLimitPropertyValue(prop v1.RateLimitProperty) string {
	res, err := yaml.Marshal(&prop)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml")
	}
	return string(res)
}

func buildConfigMapDataFileName(name string) string {
	return name + ".yaml"
}
