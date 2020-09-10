package configmapwatcher

import (
	"context"
	"github.com/stretchr/testify/assert"
	istio_v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"ratelimit-operator/pkg/apis"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/common"
	"ratelimit-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func Test_Reconcile_NotFoundRateLimiter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NotFoundRateLimiter)", func(t *testing.T) {
		r := buildEmptyReconciler()

		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		result, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(result)
	})
}

func Test_Reconcile_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		request := buildReconcileRequest(rateLimiter)
		result, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(result)
	})
}

func Test_ReconcileConfigMap_NotFoundConfigMap(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile ConfigMap (NotFoundConfigMap)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		result, err := r.reconcileConfigMap(context.TODO(), rateLimiter)

		a.Nil(err)
		a.NotNil(result)
	})
}

func Test_ReconcileConfigMap_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile ConfigMap (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		configMap := &corev1.ConfigMap{}
		configMap.Name = rateLimiter.Name
		configMap.Namespace = rateLimiter.Namespace
		r.client.Create(context.TODO(), configMap)

		result, err := r.reconcileConfigMap(context.TODO(), rateLimiter)

		a.Nil(err)
		a.NotNil(result)
	})
}

func Test_GetRateLimiterLists(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("get RateLimiterLists", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		r := buildReconciler(rateLimiter)

		list, err := r.getRateLimiterLists(context.TODO(), rateLimiter)

		a.Nil(err)
		a.NotNil(list)
	})
}

func Test_UpdateConfigMap_EmptyList(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (EmptyList)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		configMap := &corev1.ConfigMap{}
		list := &v1.RateLimiterConfigList{}
		r := buildReconciler(rateLimiter)

		r.updateConfigMap(context.TODO(), configMap, list)

		a.Nil(configMap.Data)
	})
}

func Test_UpdateConfigMap_NotEmptyList(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("update ConfigMap (NotEmptyList)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		rateLimiterConfig := buildRateLimiterConfig(rateLimiter)
		configMap := &corev1.ConfigMap{}
		list := &v1.RateLimiterConfigList{}
		list.Items = []v1.RateLimiterConfig{*rateLimiterConfig}
		fileName := common.BuildConfigMapDataFileName(rateLimiterConfig.Name)
		r := buildReconciler(rateLimiter)

		r.updateConfigMap(context.TODO(), configMap, list)

		a.NotNil(configMap.Data)
		_, found := configMap.Data[fileName]
		a.True(found)
	})
}

func buildRateLimiter() *v1.RateLimiter {
	logLevel := v1.INFO
	size := int32(1)

	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
		Spec: v1.RateLimiterSpec{
			LogLevel: &logLevel,
			Size:     &size,
		},
		Status: v1.RateLimiterStatus{},
	}
}

func buildRateLimiterConfig(rl *v1.RateLimiter) *v1.RateLimiterConfig {
	host := utils.BuildRandomString(3)
	failureModeDeny := true
	rateLimitRequestTimeout := "0.25s"
	return &v1.RateLimiterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: rl.Namespace,
		},
		Spec: v1.RateLimiterConfigSpec{
			ApplyTo:     v1.GATEWAY,
			Host:        &host,
			Port:        int32(utils.BuildRandomInt(2)),
			RateLimiter: rl.Name,
			Descriptors: []v1.Descriptor{{
				Key: utils.BuildRandomString(3),
			}},
			RateLimitRequestTimeout: &rateLimitRequestTimeout,
			FailureModeDeny:         &failureModeDeny,
		},
	}
}

func buildReconciler(rateLimiter *v1.RateLimiter) *ReconcileConfigMapWatcher {
	scheme := buildScheme()
	objects := []runtime.Object{rateLimiter}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileConfigMapWatcher{client: client, scheme: scheme}
}

func buildEmptyReconciler() *ReconcileConfigMapWatcher {
	scheme := buildScheme()
	objects := []runtime.Object{}
	client := fake.NewFakeClientWithScheme(scheme, objects...)
	return &ReconcileConfigMapWatcher{client: client, scheme: scheme}
}

func buildScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = apis.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apiextensionsv1beta1.AddToScheme(scheme)
	_ = istio_v1alpha3.AddToScheme(scheme)
	return scheme
}

func buildReconcileRequest(rateLimiter *v1.RateLimiter) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      rateLimiter.Name,
			Namespace: rateLimiter.Namespace,
		},
	}
}
