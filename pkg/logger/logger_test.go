package logger_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Logger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[logger] test suite")
}
