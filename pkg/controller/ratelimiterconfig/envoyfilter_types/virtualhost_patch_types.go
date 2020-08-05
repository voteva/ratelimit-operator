package envoyfilter_types

type Header struct {
	ExactMatch string `json:"exact_match" yaml:"exact_match"`
	Name       string `json:"name" yaml:"name"`
}

type HeaderValueMatch struct {
	DescriptorValue string   `json:"descriptor_value" yaml:"descriptor_value"`
	ExpectMatch     bool     `json:"expect_match" yaml:"expect_match"`
	Headers         []Header `json:"headers" yaml:"headers"`
}

type RequestHeader struct {
	DescriptorKey string `json:"descriptor_key" yaml:"descriptor_key"`
	HeaderName    string `json:"header_name" yaml:"header_name"`
}

type Action struct {
	RequestHeaders   *RequestHeader    `json:"request_headers,omitempty" yaml:"request_headers,omitempty"`
	HeaderValueMatch *HeaderValueMatch `json:"header_value_match,omitempty" yaml:"header_value_match,omitempty"`
}

type RateLimit struct {
	Actions []Action `json:"actions" yaml:"actions"`
}

type VirtualHostPatchValues struct {
	RateLimits []RateLimit `json:"rate_limits" yaml:"rate_limits"`
}
