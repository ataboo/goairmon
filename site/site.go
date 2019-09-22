package site

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/services/flash"
	"goairmon/business/services/identity"
	"goairmon/business/services/poll"
	"goairmon/business/services/provider"
	"goairmon/business/services/viewloader"
	"goairmon/site/controllers"
	"goairmon/site/helper"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echomiddleware "github.com/labstack/echo/middleware"
)

func EnvSiteConfig() *Config {
	return &Config{
		Address:               helper.MustGetEnv("SERVER_ADDRESS"),
		AppCookieKey:          helper.MustGetEnv("APP_COOKIE_KEY"),
		CookieStoreEncryption: helper.MustGetEnv("COOKIE_STORE_ENCRYPTION"),
		StoragePath:           helper.MustGetEnv("STORAGE_PATH"),
		SensorPointCount:      helper.MustGetEnvInt("SENSOR_POINT_COUNT"),
	}
}

func NewSite(cfg *Config) *Site {
	identityService := identity.NewIdentityService(&identity.IdentityConfig{
		CookieStoreKeySession:    cfg.AppCookieKey,
		CookieStoreEncryptionKey: cfg.CookieStoreEncryption,
	})

	site := Site{
		echoServer:      echo.New(),
		identityService: identityService,
		cfg:             cfg,
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
	SensorPointCount      int
	EncodeReadible        bool
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
	dbContext := context.NewMemDbContext(&context.MemDbConfig{
		StoragePath:      cfg.StoragePath,
		SensorPointCount: cfg.SensorPointCount,
		EncodeReadible:   cfg.EncodeReadible,
		Logger:           s.echoServer.Logger,
	})

	pollCfg := &poll.Config{
		PollDelayMillis: 60 * 1000,
		Logger:          s.echoServer.Logger,
	}
	poll := poll.NewPollService(pollCfg, dbContext)
	if err := poll.Start(); err != nil {
		s.echoServer.Logger.Info("failed to start sensor poll", err.Error())
	}

	provider.Register(viewloader.CtxKey, &viewloader.ViewLoader{})
	provider.Register(helper.CtxFlashServiceKey, flashService)
	provider.Register(helper.CtxDbContext, dbContext)
	provider.Register(helper.CtxSensorPoll, poll)

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
	s.echoServer.Static("/static", "resources/assets")
	s.echoServer.File("favicon.ico", "resources/assets/imgs/favicon.ico")

	controllers.HomeController(s.echoServer, s.identityService)
	controllers.AuthController(s.echoServer, s.identityService)
}
