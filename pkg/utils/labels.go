package utils

func LabelsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}

func SelectorsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}

func AnnotationSidecarIstio() map[string]string {
	return map[string]string{"sidecar.istio.io/inject": "true"}
}
