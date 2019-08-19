package site

import (
	"fmt"
	"goairmon/site/context"
	"goairmon/site/controllers"
	"goairmon/site/services/flash"
	"goairmon/site/services/identity"
	"goairmon/site/services/provider"
	"goairmon/site/services/viewloader"

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
	flashService := &flash.FlashService{}

	s.identity.RegisterWithProvider(provider)
	provider.Register(viewloader.CtxKey, &viewloader.ViewLoader{})
	provider.Register(context.CtxFlashServiceKey, flashService)

	s.echoServer.Use(echomiddleware.Logger())
	// s.echoServer.Use(echomiddleware.Recover())
	s.echoServer.Use(provider.BindServices())
	s.echoServer.Use(s.identity.LoadCurrentSession())
	s.echoServer.Use(flashService.PopToContext())
}

func (s *Site) bindActions() {
	s.echoServer.Static("/static", "assets")

	controllers.HomeController(s.echoServer, s.identity)
	controllers.AuthController(s.echoServer, s.identity)
}
