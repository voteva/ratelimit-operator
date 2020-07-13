package ratelimiter

import (
	"context"

	operatorsv1alpha1 "ratelimit-operator/pkg/apis/operators/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	// Watch for changes to secondary resource Pods and requeue the owner RateLimiter
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
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

	// Fetch the RateLimiter instance
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

	// Check Deployment
	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.DeploymentForRateLimiter(instance)
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

	// Check ConfigMap
	foundConfigMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new config map
		cm := r.ConfigMapForRateLimiter(instance)
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ConfigMap")
		return reconcile.Result{}, err
	}

	// Check Service
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new config map
		cm := r.ServiceForRateLimiter(instance)
		reqLogger.Info("Creating a new Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
		err = r.client.Create(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cm.Namespace, "Service.Name", cm.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	// TODO check VirtualService

	// TODO check EnvoyFilter

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
