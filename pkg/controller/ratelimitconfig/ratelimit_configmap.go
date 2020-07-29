package ratelimitconfig

import (
	"context"
	"github.com/ghodss/yaml"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimitConfig) reconcileConfigMap(ctx context.Context, instance *v1.RateLimitConfig) (reconcile.Result, error) {
	data := r.configMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	fileName := buildFileName(instance.Name)

	for key, value := range data {
		props := r.unmarshalRateLimitPropertyValue(value)

		if props.Domain == instance.Spec.RateLimitProperty.Domain && key != fileName {
			log.Error(nil, "Failed to add new rate limit configuration. Config already exists with domain "+props.Domain)
			return reconcile.Result{}, nil
		}
	}

	data[fileName] = r.buildRateLimitPropertyValue(instance)

	r.configMap.Data = data

	err := r.client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to update Config Map")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimitConfig) buildRateLimitPropertyValue(instance *v1.RateLimitConfig) string {
	res, err := yaml.Marshal(&instance.Spec.RateLimitProperty)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml")
	}
	return string(res)
}

func (r *ReconcileRateLimitConfig) unmarshalRateLimitPropertyValue(data string) v1.RateLimitProperty {
	props := v1.RateLimitProperty{}
	err := yaml.Unmarshal([]byte(data), &props)
	if err != nil {
		log.Error(err, "Failed to convert yaml to RateLimitProperty")
	}
	return props
}

func buildFileName(name string) string {
	return name + ".yaml"
}
