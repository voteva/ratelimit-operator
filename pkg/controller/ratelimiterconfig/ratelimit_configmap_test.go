package ratelimiterconfig

import (
	"github.com/stretchr/testify/assert"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildConfigMapDataFileName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ConfigMap.Data file name", func(t *testing.T) {
		fileName := utils.BuildRandomString(3)
		a.Equal(fileName + ".yaml", buildConfigMapDataFileName(fileName))
	})
}
