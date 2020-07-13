package ratelimiter

func LabelsForRateLimiter(name string) map[string]string {
	return map[string]string{"app": "ratelimiter", "ratelimiter_cr": name}
}
