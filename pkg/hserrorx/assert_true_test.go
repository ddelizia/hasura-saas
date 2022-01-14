package hserrorx_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/hserrorx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("", func() {

	logrus.SetOutput(ioutil.Discard)

	It("should not return error when condition in true", func() {
		// Given
		simpleSlice := []string{"one value"}

		// When
		err := hserrorx.AssertTrue(len(simpleSlice) == 1, hserrorx.Fields{}, "some message", nil)

		// Then
		Expect(err).Should(BeNil())
	})

	It("should return error when condition in false", func() {
		// Given
		simpleSlice := []string{"one value"}

		// When
		err := hserrorx.AssertTrue(len(simpleSlice) == 2, hserrorx.Fields{}, "some message", nil)

		// Then
		Expect(err).ShouldNot(BeNil())
	})

})
