// Package controller provides implementation of http controllers.
package controller

import (
	"net/http"

	gzEcho "github.com/gozix/echo/v2"
	gzzap "github.com/gozix/zap/v2"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"
)

// Type checking.
var _ gzEcho.Controller = (*systemController)(nil)

type systemController struct {
	log *zap.Logger
}

// DefMovieControllerName is definition name.
const DefMovieControllerName = "echo.controller.system"

// DefSystemController is provider definition getter.
func DefSystemController() di.Def {
	return di.Def{
		Name: DefMovieControllerName,
		Tags: []di.Tag{{
			Name: gzEcho.TagController,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &systemController{
				log: ctn.Get(gzzap.BundleName).(*zap.Logger).Named(DefMovieControllerName),
			}, nil
		},
	}
}

// Serve init routes
func (c *systemController) Serve(e *echo.Echo) {
	e.GET("_health", c.health)
}

// health action
func (c *systemController) health(ctx echo.Context) (err error) {

	return ctx.JSON(http.StatusOK, struct{}{})
}