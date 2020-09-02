package ratelimiterconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRatelimiterconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ratelimiterconfig Suite")
}
