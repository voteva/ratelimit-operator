package utils

func DefaultIfEmpty(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func DefaultIfAbsent(value *int32, defaultValue int32) int32 {
	if value == nil {
		return defaultValue
	}
	return *value
}
