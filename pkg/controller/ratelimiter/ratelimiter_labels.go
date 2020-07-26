package ratelimiter

func LabelsForRateLimiter(name string) map[string]string {
	return map[string]string{"app": name, "ratelimiter_cr": "ratelimiter"}
}

func LabelsForRedis(name string) map[string]string {
	return map[string]string{"app": name + "-redis", "ratelimiter_cr": "ratelimiter-redis"}
}

func SelectorsForRateLimiter(name string) map[string]string {
	return map[string]string{"app": name}
}

func SelectorsForRedis(name string) map[string]string {
	return map[string]string{"app": name + "-redis"}
}
