package utils

import "strconv"

func LabelsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}

func SelectorsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}

func AnnotationSidecarIstio() map[string]string {
	return map[string]string{"sidecar.istio.io/inject": "true"}
}

func AnnotationMetricsIstio(port uint16) map[string]string {
	return map[string]string{
		"prometheus_io_port":   strconv.Itoa(int(port)),
		"prometheus_io_scheme": "http",
		"prometheus_io_scrape": "true",
	}
}
