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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"
)

var fakeError = errors.New("fake error")

const (
	apiRoot    = "/api/"
	apiMessage = apiRoot + "message/"
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
	s.app = mocks.NewApplicationInterface(s.T())
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

func (s *HttpServerTestSuite) TestHealthCheck() {
	s.tester.GET("/").WithHeader(echo.HeaderContentType, echo.MIMETextPlain).
		Expect().
		Status(http.StatusOK).
		Text().Contains("Ok")
}

func (s *HttpServerTestSuite) TestWebhook() {
	service := gofakeit.Word()
	request := gofakeit.Sentence(3)
	s.Run("happy case", func() {
		s.app.EXPECT().ProcessRequest(mock.Anything, s.url, service, request).Return(nil).Once()
		s.tester.POST(apiRoot + s.token + "/" + service).
			WithText(request).
			Expect().
			Status(http.StatusOK).NoContent()
	})

	s.Run("wrong token", func() {
		token := gofakeit.UUID()
		s.tester.POST(apiRoot + token + "/" + service).
			WithText("test").
			Expect().
			Status(http.StatusUnauthorized)
	})

	s.Run("error in app", func() {
		s.app.EXPECT().ProcessRequest(mock.Anything, s.url, service, request).Return(fakeError).Once()
		s.tester.POST(apiRoot+s.token+"/"+service).
			WithText(request).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func (s *HttpServerTestSuite) TestShowMessage() {
	id := gofakeit.UUID()
	msg := gofakeit.Sentence(3)
	createdAt := gofakeit.Date()

	s.Run("happy case", func() {
		s.app.EXPECT().GetMessage(mock.Anything, id).Return(msg, createdAt, nil).Once()
		s.tester.GET(apiMessage+id).
			Expect().
			Status(http.StatusOK).JSON().Object().HasValue("message", msg).HasValue("created_at", createdAt).
			HasValue("id", id)
	})

	s.Run("invalid id", func() {
		invalidId := gofakeit.Word()
		s.tester.GET(apiMessage + invalidId).
			Expect().
			Status(http.StatusBadRequest).JSON().Object().ContainsKey("message")
	})

	s.Run("not found", func() {
		s.app.EXPECT().GetMessage(mock.Anything, id).Return("", time.Time{}, domain.ErrorNotFound).Once()
		s.tester.GET(apiMessage+id).
			Expect().
			Status(http.StatusNotFound).JSON().Object().HasValue("message", "not found")
	})

	s.Run("error in app", func() {
		s.app.EXPECT().GetMessage(mock.Anything, id).Return("", time.Time{}, fakeError).Once()
		s.tester.GET(apiMessage+id).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(HttpServerTestSuite))
}
