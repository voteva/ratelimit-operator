package ratelimiter

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"

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

var log = logf.Log.WithName("controller_ratelimiter")

// Add creates a new RateLimiter Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRateLimiter{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ratelimiter-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource RateLimiter
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.RateLimiter{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner RateLimiter
	log.Info("Watch for changes to appsv1.Deployment")
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	log.Info("Watch for changes to corev1.Service")
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	log.Info("Watch for changes to corev1.ConfigMap")
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.RateLimiter{},
	})
	if err != nil {
		return err
	}

	log.Info("Watch for changes to v1alpha3.VirtualService")
	err = c.Watch(&source.Kind{Type: &v1alpha3.VirtualService{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &operatorsv1alpha1.RateLimiter{},
		})
	if err != nil {
		return err
	}

	log.Info("Watch for changes to v1alpha3.EnvoyFilter")
	err = c.Watch(&source.Kind{Type: &v1alpha3.EnvoyFilter{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &operatorsv1alpha1.RateLimiter{},
		})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRateLimiter{}

type ReconcileRateLimiter struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileRateLimiter) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling RateLimiter")

	instance := &operatorsv1alpha1.RateLimiter{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if _, err := r.reconcileDeployment(request, instance); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := r.reconcileConfigMap(request, instance); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := r.reconcileService(request, instance); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := r.reconcileVirtualService(request, instance); err != nil {
		return reconcile.Result{}, err
	}

	if _, err := r.reconcileEnvoyFilter(request, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
