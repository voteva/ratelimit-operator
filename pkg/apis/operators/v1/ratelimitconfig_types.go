package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type RateLimitConfigSpec struct {
	RateLimiter       string            `json:"rateLimiter"`
	RateLimitProperty RateLimitProperty `json:"rateLimitProperty,omitempty"`
	FailureModeDeny   bool              `json:"failureModeDeny,omitempty"`
	VirtualHostName   string            `json:"virtualHostName,omitempty"`
}

type RateLimitConfigStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RateLimitConfig is the Schema for the ratelimitconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ratelimitconfigs,scope=Namespaced
type RateLimitConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateLimitConfigSpec   `json:"spec,omitempty"`
	Status RateLimitConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RateLimitConfigList contains a list of RateLimitConfig
type RateLimitConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimitConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RateLimitConfig{}, &RateLimitConfigList{})
}
