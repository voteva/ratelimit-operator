package envoyfilter_types

type RequestHeader struct {
	DescriptorKey string `json:"descriptor_key" yaml:"descriptor_key"`
	HeaderName    string `json:"header_name" yaml:"header_name"`
}

type Action struct {
	RequestHeaders RequestHeader `json:"request_headers" yaml:"request_headers"`
}

type RateLimit struct {
	Actions []Action `json:"actions" yaml:"actions"`
}

type VirtualHostPatchValues struct {
	RateLimits []RateLimit `json:"rate_limits" yaml:"rate_limits"`
}
