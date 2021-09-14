package gqlreq_test

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	machinebox "github.com/machinebox/graphql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

// graphQlClientMock
type graphQlClientMock struct {
	mock.Mock
}

func (m *graphQlClientMock) Run(c context.Context, req *machinebox.Request, res interface{}) error {
	args := m.Called(c, req, res)
	return args.Error(0)
}

var _ = Describe("ExecuterImpl.Execute()", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		s       *gqlreq.ExecuterImpl
		gqlMock *graphQlClientMock
	)

	BeforeEach(func() {
		gqlMock = &graphQlClientMock{}
		s = &gqlreq.ExecuterImpl{
			Client: gqlMock,
		}
	})

	mockRun := func(err error, arguments ...interface{}) {
		gqlMock.On("Run", arguments...).Return(err)
	}

	Context("Anything", func() {
		It("should get error when request is on error", func() {
			// Given
			mockRun(errors.New("some error"), mock.Anything, mock.Anything, mock.Anything)

			// When
			err := s.Execute(context.Background(), "somequery", nil, nil, false, &gqlreq.InsertBaseResponse{})

			// Then
			Expect(len(gqlMock.Calls)).To(Equal(1))
			Expect(err).NotTo(BeNil())
		})

		It("should set the data in the response struct passed as parameter", func() {
			// Given
			data := &gqlreq.InsertBaseResponse{}

			gqlMock.On("Run", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					data := args.Get(2).(*gqlreq.InsertBaseResponse)
					data.AffectedRows = 1
					data.Returning = []gqlreq.InsertBaseResponseID{{ID: "some id"}}
				}).Return(nil)

			// When
			err := s.Execute(context.Background(), "somequery", nil, nil, false, data)

			// Then
			Expect(len(gqlMock.Calls)).To(Equal(1))
			Expect(data).To(Equal(&gqlreq.InsertBaseResponse{
				AffectedRows: 1,
				Returning:    []gqlreq.InsertBaseResponseID{{ID: "some id"}},
			}))
			Expect(err).To(BeNil())
		})

		It("should set the admin header in the request", func() {
			// Given
			data := &gqlreq.InsertBaseResponse{}

			gqlMock.On("Run", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					data := args.Get(2).(*gqlreq.InsertBaseResponse)
					data.AffectedRows = 1
					data.Returning = []gqlreq.InsertBaseResponseID{{ID: "some id"}}
				}).Return(nil)

			// When
			err := s.Execute(context.Background(), "somequery", nil, nil, true, data)

			// Then
			mbReq := gqlMock.Calls[0].Arguments[1].(*machinebox.Request)
			Expect(len(gqlMock.Calls)).To(Equal(1))
			Expect(mbReq.Header.Get("X-Hasura-Admin-Secret")).To(Equal(gqlreq.ConfigAdminSecret()))
			Expect(err).To(BeNil())
		})
	})
})
