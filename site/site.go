package site

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/services/flash"
	"goairmon/business/services/identity"
	"goairmon/business/services/provider"
	"goairmon/business/services/viewloader"
	"goairmon/site/controllers"
	"goairmon/site/helper"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echomiddleware "github.com/labstack/echo/middleware"
)

func NewSite(cfg *Config) *Site {
	identityService := identity.NewIdentityService(&identity.IdentityConfig{
		CookieStoreKeySession:    cfg.AppCookieKey,
		CookieStoreEncryptionKey: cfg.CookieStoreEncryption,
	})

	site := Site{
		echoServer:      echo.New(),
		identityService: identityService,
	}

	site.bindGlobalMiddleware(cfg)
	site.bindActions()

	return &site
}

type Site struct {
	echoServer      *echo.Echo
	identityService *identity.IdentityService
	cfg             *Config
}

type Config struct {
	AppCookieKey          string
	CookieStoreEncryption string
	Address               string
	StoragePath           string
}

func (s *Site) Start() {
	go func() {
		s.echoServer.Logger.Fatal(s.echoServer.Start(s.cfg.Address))
	}()
}

func (s *Site) Cleanup() error {
	fmt.Print("Running cleanup!\n")
	return s.echoServer.Close()
}

func (s *Site) bindGlobalMiddleware(cfg *Config) {
	provider := provider.NewServiceProvider()
	flashService := &flash.FlashService{}
	dbContext := context.NewMemDbContext(&context.MemDbConfig{StoragePath: cfg.StoragePath})

	s.identityService.RegisterWithProvider(provider)
	provider.Register(viewloader.CtxKey, &viewloader.ViewLoader{})
	provider.Register(helper.CtxFlashServiceKey, flashService)
	provider.Register(helper.CtxDbContext, dbContext)

	s.echoServer.Use(echomiddleware.Logger())
	// s.echoServer.Use(echomiddleware.Recover())
	s.echoServer.Use(provider.BindServices())
	s.echoServer.Use(s.identityService.LoadCurrentSession())
	s.echoServer.Use(flashService.PopToContext())
	s.echoServer.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_csrf-token",
	}))

}

func (s *Site) bindActions() {
	s.echoServer.Static("/static", "site/assets")
	s.echoServer.File("favicon.ico", "site/assets/imgs/favicon.ico")

	controllers.HomeController(s.echoServer, s.identityService)
	controllers.AuthController(s.echoServer, s.identityService)
}
