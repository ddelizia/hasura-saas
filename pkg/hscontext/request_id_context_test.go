package hscontext_test

import (
	"context"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("", func() {

	logrus.SetOutput(ioutil.Discard)

	It("should get the request id once stored in context", func() {
		// Given
		ctx := context.Background()
		ctx = hscontext.WithRequestIDValue(ctx, "request id")

		// When
		requestId := hscontext.RequestIDValue(ctx)

		// Then

		Expect(requestId).To(Equal("request id"))
	})

})
