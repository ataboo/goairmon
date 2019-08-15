package site

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewSite(cfg *Config) *Site {
	site := Site{config: cfg}

	return &site
}

type Config struct {
	SessionKey   string
	DbConnection string
	Port         int
}

type Site struct {
	config     *Config
	echoServer *echo.Echo
}

func (s *Site) Start() {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))


	fmt.Printf("Hello world!")

	s.echoServer = echo.New()
	s.bindGlobalMiddleware()
	s.bindActions()

	go func() {
		s.echoServer.Logger.Fatal(s.echoServer.Start(":" + strconv.Itoa(s.config.Port)))
	}()
}

func (s *Site) bindGlobalMiddleware() {
	s.echoServer.Use(middleware.Logger())
	s.echoServer.Use(middleware.Recover())
}

func (s *Site) bindActions() {
	s.echoServer.Static("/static", "assets")

	s.echoServer.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!\n")
	})
}

func (s *Site) Cleanup() error {
	fmt.Print("Running cleanup!\n")
	return s.echoServer.Close()
}
