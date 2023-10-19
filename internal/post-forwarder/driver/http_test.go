package driver

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
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
	url    string
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
	time.Sleep(time.Second)
	s.url = "http://localhost:" + strconv.Itoa(port)
	s.tester = httpexpect.Default(s.T(), s.url)
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
		s.app.On("ProcessRequest", mock.Anything, s.url, service, request).Return(nil).Once()
		s.tester.POST("/api/" + s.token + "/" + service).
			WithText(request).
			Expect().
			Status(http.StatusOK).NoContent()
	})

	s.Run("wrong token", func() {
		token := gofakeit.UUID()
		s.tester.POST("/api/" + token + "/" + service).
			WithText("test").
			Expect().
			Status(http.StatusUnauthorized)
	})

	s.Run("error in app", func() {
		s.app.On("ProcessRequest", mock.Anything, s.url, service, request).Return(errors.New("fake error")).Once()
		s.tester.POST("/api/"+s.token+"/"+service).
			WithText(request).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", "fake error")
	})
}

func (s *HttpServerTestSuite) TestShowMessage() {
	id := gofakeit.UUID()
	msg := gofakeit.SentenceSimple()
	createdAt := gofakeit.Date()

	s.Run("happy case", func() {
		s.app.On("ShowMessage", mock.Anything, id).Return(msg, createdAt, nil).Once()
		s.tester.GET("/api/message/"+id).
			Expect().
			Status(http.StatusOK).JSON().Object().HasValue("message", msg).HasValue("created_at", createdAt).
			HasValue("id", id)
	})

	s.Run("invalid id", func() {
		invalidId := gofakeit.Word()
		s.tester.GET("/api/message/" + invalidId).
			Expect().
			Status(http.StatusBadRequest).JSON().Object().ContainsKey("message")
	})

	s.Run("not found", func() {
		s.app.On("ShowMessage", mock.Anything, id).Return("", time.Time{}, domain.ErrorNotFound).Once()
		s.tester.GET("/api/message/"+id).
			Expect().
			Status(http.StatusNotFound).JSON().Object().HasValue("message", "not found")
	})

	s.Run("error in app", func() {
		s.app.On("ShowMessage", mock.Anything, id).Return("", time.Time{}, errors.New("fake error")).Once()
		s.tester.GET("/api/message/"+id).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", "fake error")
	})
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(HttpServerTestSuite))
}
