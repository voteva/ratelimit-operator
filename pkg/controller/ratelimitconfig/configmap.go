package ratelimitconfig

import (
	"context"
	"github.com/ghodss/yaml"
	operatorsv1 "ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimitConfig) reconcileConfigMap(ctx context.Context, instance *operatorsv1.RateLimitConfig) (reconcile.Result, error) {
	data := r.configMap.Data
	data[instance.Name + ".yaml"] = r.buildRateLimitPropertyValue(instance)

	r.configMap.Data = data

	err := r.client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to update Config Map")
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, nil
}

func (r *ReconcileRateLimitConfig) buildRateLimitPropertyValue(instance *operatorsv1.RateLimitConfig) string {
	res, err := yaml.Marshal(&instance.Spec.RateLimitProperty)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml")
	}
	return string(res)
}
