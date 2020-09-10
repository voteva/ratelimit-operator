package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplyTo string

const (
	GATEWAY          ApplyTo = "GATEWAY"
	SIDECAR_INBOUND  ApplyTo = "SIDECAR_INBOUND"
	SIDECAR_OUTBOUND ApplyTo = "SIDECAR_OUTBOUND"
)

type RateLimits struct {
	DisableKey string   `json:"disable_key,omitempty" yaml:"disable_key,omitempty"`
	Actions    []Action `json:"actions" yaml:"actions"`
}

type WorkloadSelector struct {
	Labels map[string]string `json:"labels"`
}

type RateLimiterConfigSpec struct {
	// +kubebuilder:validation:Enum={GATEWAY,SIDECAR_INBOUND,SIDECAR_OUTBOUND}
	ApplyTo ApplyTo `json:"applyTo"`
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:MinLength=1
	Host *string `json:"host,omitempty"`
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=0
	Port int32 `json:"port"`
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:MinLength=1
	RateLimiter string `json:"rateLimiter"`
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	RateLimitRequestTimeout *string          `json:"rateLimitRequestTimeout,omitempty"`
	FailureModeDeny         *bool            `json:"failureModeDeny,omitempty"`
	WorkloadSelector        WorkloadSelector `json:"workloadSelector"`
	Descriptors             []Descriptor     `json:"descriptors,omitempty" yaml:"descriptors,omitempty"`
	RateLimits              []RateLimits     `json:"rateLimits,omitempty" yaml:"rateLimits,omitempty"`
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
