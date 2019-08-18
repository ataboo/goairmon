package middleware

import "github.com/labstack/echo"

type ServiceProvider struct {
	bindings map[string]interface{}
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{
		bindings: make(map[string]interface{}),
	}
}

func (p *ServiceProvider) Register(key string, service interface{}) {
	p.bindings[key] = service
}

func (p *ServiceProvider) BindServices() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for k, v := range p.bindings {
				c.Set(k, v)
			}

			return next(c)
		}
	}
}
