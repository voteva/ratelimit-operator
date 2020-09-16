package apis

import (
	v1 "github.com/voteva/ratelimit-operator/pkg/apis/operators/v1"
)

func init() {
	// Register the types with the scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
