package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RateLimit struct {
	Unit            string `json:"unit"`
	RequestsPerUnit int32  `json:"requests_per_unit"`
}

type Descriptor struct {
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	RateLimit RateLimit `json:"rate_limit,omitempty"`
	Descriptors []Descriptor `json:"descriptors,omitempty"`
}

type RateLimitProperty struct {
	Domain      string       `json:"domain"`
	Descriptors []Descriptor `json:"descriptors"`
}

// RateLimiterSpec defines the desired state of RateLimiter
type RateLimiterSpec struct {
	Size              int32             `json:"size"`
	ServicePort       int32             `json:"servicePort"`
	RateLimitProperty RateLimitProperty `json:"rateLimitProperty"`
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
