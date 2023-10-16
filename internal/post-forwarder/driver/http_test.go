package driver

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"
)

type HttpServerTestSuite struct {
	suite.Suite
	e      *echo.Echo
	tester *httpexpect.Expect
	token  string
	app    *mocks.ApplicationInterface
}

func (s *HttpServerTestSuite) SetupSuite() {
	s.token = gofakeit.UUID()
	s.app = new(mocks.ApplicationInterface)
	s.e = echo.New()
	RegisterHandlers(s.e, NewHTTPServer(&config.Config{
		AuthToken: s.token,
	}, s.app))
	port, err := freeport.GetFreePort()
	s.Require().NoError(err)
	go func() {
		err := s.e.Start(":" + strconv.Itoa(port))
		s.Require().True(err == nil || errors.Is(err, http.ErrServerClosed))
	}()

	s.tester = httpexpect.Default(s.T(), "http://localhost:"+strconv.Itoa(port))
}

func (s *HttpServerTestSuite) TearDownSuite() {
	err := s.e.Shutdown(context.Background())
	s.NoError(err)
}

func (s *HttpServerTestSuite) TearDownTest() {
	s.app.AssertExpectations(s.T())
}

func (s *HttpServerTestSuite) TestHealthCheck() {
	s.tester.GET("/").WithHeader(echo.HeaderContentType, echo.MIMETextPlain).
		Expect().
		Status(http.StatusOK).
		Text().Contains("Ok")
}

func (s *HttpServerTestSuite) TestWebhook() {
	service := gofakeit.Word()
	request := gofakeit.SentenceSimple()
	s.Run("happy case", func() {
		s.app.On("ProcessRequest", mock.Anything, service, request).Return(nil).Once()
		s.tester.POST("/api/"+s.token+"/"+service).WithHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
			WithText(request).
			Expect().
			Status(http.StatusOK).NoContent()
	})

	s.Run("wrong token", func() {
		token := gofakeit.UUID()
		s.tester.POST("/api/"+token+"/"+service).WithHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
			WithText("test").
			Expect().
			Status(http.StatusUnauthorized)
	})

	s.Run("error in app", func() {
		s.app.On("ProcessRequest", mock.Anything, service, request).Return(errors.New("fake error")).Once()
		s.tester.POST("/api/"+s.token+"/"+service).WithHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
			WithText(request).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", "fake error")
	})
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(HttpServerTestSuite))
}
