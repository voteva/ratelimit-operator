package utils

func Merge(map1 map[string]string, map2 map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}
