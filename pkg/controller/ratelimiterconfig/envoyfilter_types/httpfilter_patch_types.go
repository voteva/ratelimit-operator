package envoyfilter_types

type EnvoyGrpc struct {
	ClusterName string `json:"cluster_name" yaml:"cluster_name"`
}

type GrpcService struct {
	EnvoyGrpc EnvoyGrpc `json:"envoy_grpc" yaml:"envoy_grpc"`
	Timeout   string    `json:"timeout" yaml:"timeout"`
}

type RateLimitService struct {
	GrpcService GrpcService `json:"grpc_service" yaml:"grpc_service"`
}

type Config struct {
	Domain           string           `json:"domain" yaml:"domain"`
	FailureModeDeny  bool             `json:"failure_mode_deny" yaml:"failure_mode_deny"`
	RateLimitService RateLimitService `json:"rate_limit_service" yaml:"rate_limit_service"`
}

type HttpFilterPatchValues struct {
	Config Config `json:"config" yaml:"config"`
	Name   string `json:"name" yaml:"name"`
}
