package ratelimiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func Test_Reconcile_NotFoundRateLimiter(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (NotFoundRateLimiter)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
	})
}

func Test_Reconcile_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("reconcile (Success)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		request := buildReconcileRequest(rateLimiter)
		r := buildEmptyReconciler()

		errCreate := r.client.Create(context.Background(), rateLimiter)
		a.Nil(errCreate)

		reconcileResult, err := r.Reconcile(request)

		a.Nil(err)
		a.NotNil(reconcileResult)
	})
}

func buildReconcileRequest(rateLimiter *v1.RateLimiter) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      rateLimiter.Name,
			Namespace: rateLimiter.Namespace,
		},
	}
}
