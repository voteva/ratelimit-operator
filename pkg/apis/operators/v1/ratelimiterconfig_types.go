package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplyTo string

const (
	GATEWAY          ApplyTo = "GATEWAY"
	SIDECAR_OUTBOUND ApplyTo = "SIDECAR_OUTBOUND"
	SIDECAR_INBOUND  ApplyTo = "SIDECAR_INBOUND"
)

type RateLimit struct {
	Unit            string `json:"unit" yaml:"unit"`
	RequestsPerUnit int32  `json:"requests_per_unit" yaml:"requests_per_unit"`
}

type DescriptorInternal struct {
	Key       string    `json:"key" yaml:"key"`
	Value     string    `json:"value,omitempty" yaml:"value,omitempty"`
	RateLimit RateLimit `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
}

type Descriptor struct {
	Key         string               `json:"key" yaml:"key"`
	Value       string               `json:"value,omitempty" yaml:"value,omitempty"`
	RateLimit   RateLimit            `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
	Descriptors []DescriptorInternal `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}

type RateLimitProperty struct {
	Domain      string       `json:"domain" yaml:"domain"`
	Descriptors []Descriptor `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
}

type WorkloadSelector struct {
	Labels map[string]string `json:"labels"`
}

type RateLimiterConfigSpec struct {
	ApplyTo           ApplyTo           `json:"applyTo"`
	Port              int32             `json:"port"`
	Host              *string           `json:"host,omitempty"`
	WorkloadSelector  *WorkloadSelector `json:"workloadSelector,omitempty"`
	RateLimitProperty RateLimitProperty `json:"rateLimitProperty,omitempty"`
	FailureModeDeny   bool              `json:"failureModeDeny,omitempty"`
	RateLimiter       string            `json:"rateLimiter"`
}

type RateLimiterConfigStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RateLimiterConfig is the Schema for the ratelimiterconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ratelimiterconfigs,scope=Namespaced
type RateLimiterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateLimiterConfigSpec   `json:"spec,omitempty"`
	Status RateLimiterConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RateLimiterConfigList contains a list of RateLimiterConfig
type RateLimiterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimiterConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RateLimiterConfig{}, &RateLimiterConfigList{})
}
