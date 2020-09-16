package utils

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
	"testing"
)

func Test_IsBeingDeleted(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("IsBeingDeleted", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		result := IsBeingDeleted(rateLimiter)

		a.False(result)
	})
}

func Test_HasFinalizer_False(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("HasFinalizer (False)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		finalizer := BuildRandomString(3)
		AddFinalizer(rateLimiter, finalizer)
		RemoveFinalizer(rateLimiter, finalizer)
		result := HasFinalizer(rateLimiter, finalizer)

		a.False(result)
	})
}

func Test_HasFinalizer_True(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("HasFinalizer (True)", func(t *testing.T) {
		rateLimiter := buildRateLimiter()
		finalizer := BuildRandomString(3)
		AddFinalizer(rateLimiter, finalizer)
		result := HasFinalizer(rateLimiter, finalizer)

		a.True(result)
	})
}

func buildRateLimiter() *v1.RateLimiter {
	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BuildRandomString(3),
			Namespace: BuildRandomString(3),
		},
		Spec:   v1.RateLimiterSpec{},
		Status: v1.RateLimiterStatus{},
	}
}
