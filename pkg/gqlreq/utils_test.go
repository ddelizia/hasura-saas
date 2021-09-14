package gqlreq

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/hstype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("util", func() {

	logrus.SetOutput(ioutil.Discard)

	Context("graphql.TimestampToGraphQlTimestamWithTimeszone()", func() {
		It("should create a valid string date according to the format", func() {
			// TODO fix this test on the pipeline
			Skip("Skipping for the moment cause timezon on the pipeline")
			// When
			got := TimestampToGraphQlTimestamWithTimeszone(1621798437)

			// Then
			Expect(got).To(Equal("2021-05-23T21:33:57+02:00"))
		})
	})

	Context("graphql.InterfaceToJson()", func() {
		It("should convert map to json string", func() {
			// When
			got, err := InterfaceToJson(map[string]hstype.String{
				"key1": hstype.NewString("value1"),
				"key2": hstype.NewString("value2"),
			})

			// Then
			Expect(got).To(Equal("{\"key1\":\"value1\",\"key2\":\"value2\"}"))
			Expect(err).To(BeNil())
		})
	})

})
