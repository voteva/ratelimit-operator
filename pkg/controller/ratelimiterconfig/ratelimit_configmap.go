package ratelimiterconfig

import (
	"context"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/common"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiterConfig) updateConfigMap(ctx context.Context, instance *v1.RateLimiterConfig) (reconcile.Result, error) {
	data := r.configMap.Data
	if data == nil {
		data = make(map[string]string)
	}

	fileName := common.BuildConfigMapDataFileName(instance.Name)
	data[fileName] = common.BuildRateLimitPropertyValue(instance)

	r.configMap.Data = data

	err := r.Client.Update(ctx, r.configMap)
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

	fileName := common.BuildConfigMapDataFileName(instance.Name)
	delete(data, fileName)

	r.configMap.Data = data

	err := r.Client.Update(ctx, r.configMap)
	if err != nil {
		log.Error(err, "Failed to delete keys from Config Map for RateLimiterConfig [%s]", instance.Name)
		return err
	}
	return nil
}
