package ratelimiterconfig

import (
	"context"
	"github.com/ghodss/yaml"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiterConfig) updateConfigMap(ctx context.Context, instance *v1.RateLimiterConfig) (reconcile.Result, error) {
	data := r.configMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	fileName := buildFileName(instance.Name)

	for key, value := range data {
		props := unmarshalRateLimitPropertyValue(value)

		if props.Domain == instance.Spec.RateLimitProperty.Domain && key != fileName {
			log.Error(nil, "Failed to add new rate limit configuration. Config already exists with domain "+props.Domain)
			return reconcile.Result{}, nil
		}
	}

	data[fileName] = buildRateLimitPropertyValue(instance.Spec.RateLimitProperty)

	r.configMap.Data = data

	err := r.client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to update Config Map")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiterConfig) deleteFromConfigMap(ctx context.Context, instance *v1.RateLimiterConfig) error {
	data := r.configMap.Data
	if data == nil {
		return nil
	}

	fileName := buildFileName(instance.Name)
	delete(data, fileName)

	r.configMap.Data = data

	err := r.client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to delete keys from Config Map for RateLimiterConfig [%s]", instance.Name)
		return err
	}
	return nil
}

func buildRateLimitPropertyValue(prop v1.RateLimitProperty) string {
	res, err := yaml.Marshal(&prop)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml")
	}
	return string(res)
}

func unmarshalRateLimitPropertyValue(data string) v1.RateLimitProperty {
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
