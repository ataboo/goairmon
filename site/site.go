package site

import (
	"fmt"
	"goairmon/site/identity"
	"goairmon/site/middleware"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	echomiddleware "github.com/labstack/echo/middleware"
)

func NewSite() *Site {
	identityCfg := identity.DefaultIdentityCfg()
	identityCfg.CookieStoreKeySession = mustGetEnv("APP_COOKIE_KEY")
	identityCfg.CookieStoreEncryptionKey = mustGetEnv("COOKIE_STORE_ENCRYPTION")

	site := Site{
		echoServer: echo.New(),
		identity:   identity.NewIdentityService(identityCfg),
	}

	return &site
}

type Site struct {
	echoServer *echo.Echo
	identity   *identity.IdentityService
}

func (s *Site) Start() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("failed to load .env")
	}

	s.echoServer = echo.New()
	s.bindGlobalMiddleware()
	s.bindActions()

	go func() {
		s.echoServer.Logger.Fatal(s.echoServer.Start(":" + mustGetEnv("SERVER_PORT")))
	}()
}

func (s *Site) bindGlobalMiddleware() {
	provider := middleware.NewServiceProvider()

	s.identity.RegisterWithProvider(provider)

	s.echoServer.Use(echomiddleware.Logger())
	s.echoServer.Use(echomiddleware.Recover())
	s.echoServer.Use(provider.BindServices())
	s.identity.RegisterWithProvider(provider)
	s.echoServer.Use(s.identity.LoadCurrentSession())
}

func (s *Site) bindActions() {
	s.echoServer.Static("/static", "assets")

	s.echoServer.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})

	authorizeFail := func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "you have to log in")
	}
	secure := s.echoServer.Group("/secure", s.identity.RequireSession(authorizeFail))
	secure.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Secure World\n")
	})
	secure.GET("/login", func(c echo.Context) error {

		return c.String(http.StatusOK, "started session")
	})
	secure.GET("/logout", func(c echo.Context) error {

		return nil
	})
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Failed to load dotenv value: %s", key))
	}

	return val
}

func (s *Site) Cleanup() error {
	fmt.Print("Running cleanup!\n")
	return s.echoServer.Close()
}
