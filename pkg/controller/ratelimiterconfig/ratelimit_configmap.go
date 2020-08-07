package ratelimiterconfig

import (
	"context"
	"github.com/ghodss/yaml"
	"ratelimit-operator/pkg/apis/operators/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

func (r *ReconcileRateLimiterConfig) updateConfigMap(ctx context.Context, instance *v1.RateLimiterConfig) (reconcile.Result, error) {
	data := r.configMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	for _, patch := range instance.Spec.ConfigPatches {
		fileName := buildFileName(instance.Name, patch)

		for key, value := range data {
			props := r.unmarshalRateLimitPropertyValue(value)

			if props.Domain == patch.RateLimitProperty.Domain && key != fileName {
				log.Error(nil, "Failed to add new rate limit configuration. Config already exists with domain "+props.Domain)
				return reconcile.Result{}, nil
			}
		}

		data[fileName] = r.buildRateLimitPropertyValue(patch.RateLimitProperty)

		r.configMap.Data = data
	}

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

	for _, patch := range instance.Spec.ConfigPatches {
		fileName := buildFileName(instance.Name, patch)
		delete(data, fileName)
	}

	r.configMap.Data = data

	err := r.client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to delete keys from Config Map for RateLimiterConfig [%s]", instance.Name)
		return err
	}
	return nil
}

func (r *ReconcileRateLimiterConfig) buildRateLimitPropertyValue(prop v1.RateLimitProperty) string {
	res, err := yaml.Marshal(&prop)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) unmarshalRateLimitPropertyValue(data string) v1.RateLimitProperty {
	props := v1.RateLimitProperty{}
	err := yaml.Unmarshal([]byte(data), &props)
	if err != nil {
		log.Error(err, "Failed to convert yaml to RateLimitProperty")
	}
	return props
}

func buildFileName(name string, patch v1.ConfigPatch) string {
	return name + "-" + strings.ToLower(string(patch.ApplyTo)) + ".yaml"
}
