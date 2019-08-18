package site

import (
	"fmt"
	"goairmon/site/services/identity"
	"goairmon/site/services/provider"
	"log"
	"net/http"

	"github.com/labstack/echo"
	echomiddleware "github.com/labstack/echo/middleware"
)

func NewSite() *Site {
	site := Site{
		echoServer: echo.New(),
		identity:   identity.NewIdentityService(nil),
	}

	return &site
}

type Site struct {
	echoServer *echo.Echo
	identity   *identity.IdentityService
}

type Config struct {
	AppCookieKey          string
	CookieStoreEncryption string
	Address               string
}

func (s *Site) Start(cfg *Config) {
	s.identity.Cfg.CookieStoreKeySession = cfg.AppCookieKey
	s.identity.Cfg.CookieStoreEncryptionKey = cfg.CookieStoreEncryption

	s.echoServer = echo.New()
	s.bindGlobalMiddleware()
	s.bindActions()

	go func() {
		s.echoServer.Logger.Fatal(s.echoServer.Start(cfg.Address))
	}()
}

func (s *Site) Cleanup() error {
	fmt.Print("Running cleanup!\n")
	return s.echoServer.Close()
}

func (s *Site) bindGlobalMiddleware() {
	provider := provider.NewServiceProvider()

	s.identity.RegisterWithProvider(provider)

	s.echoServer.Use(echomiddleware.Logger())
	// s.echoServer.Use(echomiddleware.Recover())
	s.echoServer.Use(provider.BindServices())
	s.echoServer.Use(s.identity.LoadCurrentSession())
}

func (s *Site) bindActions() {
	s.echoServer.Static("/static", "assets")

	s.echoServer.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})

	//TODO: change me to post
	s.echoServer.GET("/login", func(c echo.Context) error {
		err := s.identity.StartNewSession(c)
		if err != nil {
			log.Println("failed to login:", err)
			return c.String(http.StatusUnauthorized, "failed to log in")
		}

		return c.String(http.StatusOK, "started session")
	})

	//TODO: change me to post
	s.echoServer.GET("/logout", func(c echo.Context) error {
		err := s.identity.EndSession(c)
		if err != nil {
			log.Println("failed to logout:", err)
			return c.String(http.StatusUnauthorized, "failed to log out")
		}

		return c.String(http.StatusOK, "ended session")
	})

	secure := s.echoServer.Group("/secure", s.identity.RequireSession(func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "you have to log in")
	}))
	secure.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Secure World\n")
	})
}
