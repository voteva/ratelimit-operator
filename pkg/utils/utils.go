package utils

func DefaultIfEmpty(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

func DefaultIfZero(value int32, defaultValue int32) int32 {
	if value > 0 {
		return value
	}
	return defaultValue
}
