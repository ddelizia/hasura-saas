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
		sessionVars := map[string]interface{}{
			"key": "value",
		}
		ctx = hscontext.WithActionSessionValue(ctx, sessionVars, "some value")

		// When
		s := hscontext.ActionSessionVariablesValue(ctx)
		d := hscontext.ActionDataValue(ctx)

		// Then

		Expect(s).To(Equal(sessionVars))
		Expect(d).To(Equal("some value"))
	})

})
