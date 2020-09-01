package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LabelsForApp(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	label := LabelsForApp("appname")
	a.Equal(label, map[string]string{"app": "appname"})
}

func Test_SelectorForApp(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	label := SelectorsForApp("appname")
	a.Equal(label, map[string]string{"app": "appname"})
}

func Test_AnnotationSidecarIstio(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	label := AnnotationSidecarIstio()
	a.Equal(label, map[string]string{"sidecar.istio.io/inject": "true"})
}

func Test_AnnotationMetricsIstio(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	label := AnnotationMetricsIstio(9102)
	a.Equal(label, map[string]string{
		"prometheus_io_port":   "9102",
		"prometheus_io_scheme": "http",
		"prometheus_io_scrape": "true",
	})
}
