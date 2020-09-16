package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/voteva/ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildRateLimitPropertyValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build RateLimitProperty value", func(t *testing.T) {
		rl := buildRateLimiter()
		rlc := buildRateLimiterConfig(rl)
		a.NotNil(BuildRateLimitPropertyValue(rlc))
	})
}

func Test_BuildConfigMapDataFileName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build ConfigMap.Data file name", func(t *testing.T) {
		fileName := utils.BuildRandomString(3)
		a.Equal(fileName+".yaml", BuildConfigMapDataFileName(fileName))
	})
}
