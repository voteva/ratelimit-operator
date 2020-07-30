package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Resource interface {
	metav1.Object
	runtime.Object
}

func IsBeingDeleted(obj Resource) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

func HasFinalizer(obj Resource, finalizer string) bool {
	for _, fin := range obj.GetFinalizers() {
		if fin == finalizer {
			return true
		}
	}
	return false
}

func AddFinalizer(obj Resource, finalizer string) {
	controllerutil.AddFinalizer(obj, finalizer)
}

func RemoveFinalizer(obj Resource, finalizer string) {
	controllerutil.RemoveFinalizer(obj, finalizer)
}
