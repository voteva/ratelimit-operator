package ratelimiter

func LabelsForRateLimiter(name string) map[string]string {
	return map[string]string{"app": name, "ratelimiter_cr": "ratelimiter"}
}

func SelectorsForRateLimiter(name string) map[string]string {
	return map[string]string{"app": name}
}
