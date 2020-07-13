package controller

import (
	"ratelimit-operator/pkg/controller/ratelimiter"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ratelimiter.Add)
}
