package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LogLevel string

const (
	// Log level INFO.
	INFO LogLevel = "INFO"
	// Log level DEBUG.
	DEBUG LogLevel = "DEBUG"
	// Log level ERROR.
	ERROR LogLevel = "ERROR"
)

// RateLimiterSpec defines the desired state of RateLimiter
type RateLimiterSpec struct {
	Port     *int32    `json:"port,omitempty"`
	LogLevel *LogLevel `json:"logLevel,omitempty"`
}

// RateLimiterStatus defines the observed state of RateLimiter
type RateLimiterStatus struct {
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RateLimiter is the Schema for the ratelimiters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ratelimiters,scope=Namespaced
type RateLimiter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateLimiterSpec   `json:"spec,omitempty"`
	Status RateLimiterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RateLimiterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimiter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RateLimiter{}, &RateLimiterList{})
}
