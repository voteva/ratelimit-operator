package controller

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"testing"
)

func TestNewStubClient(t *testing.T) {

	c := v1.ConfigMap{}

	eventHandler := StubEventHandler{
		EventHandler: &handler.EnqueueRequestForOwner{
			OwnerType:    &c,
			IsController: true,
		},
		Queue: NewStubQueue(),
	}

	isController := true
	rateLimiter_1 := v12.RateLimiter{
		TypeMeta: metav1.TypeMeta{
			Kind: "RateLimiter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimiter",
			Namespace: "test-namespace",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Controller: &isController,
			}},
		},
	}

	rateLimiter_2 := v12.RateLimiter{
		TypeMeta: metav1.TypeMeta{
			Kind: "RateLimiter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ratelimiter",
			Namespace: "test-namespace",
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Controller: &isController,
			}},
		},
	}

	eventHandler.OnUpdate(&rateLimiter_1, &rateLimiter_2)
}
