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
var _ gzEcho.Controller = (*movieController)(nil)

//easyjson:json
type (
	response struct{}
)
type movieController struct {
	log *zap.Logger
}

// DefMovieControllerName is definition name.
const DefMovieControllerName = "echo.controller.movies"

// DefMovieController is provider definition getter.
func DefMovieController() di.Def {
	return di.Def{
		Name: DefMovieControllerName,
		Tags: []di.Tag{{
			Name: gzEcho.TagController,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &movieController{
				log: ctn.Get(gzzap.BundleName).(*zap.Logger).Named(DefMovieControllerName),
			}, nil
		},
	}
}

// Serve init routes
func (c *movieController) Serve(e *echo.Echo) {
	e.GET("/api/v1/movies", c.CreateMovie)
}

// CreateMovie action
//   curl -X POST 'http://localhost:8888/api/v1/movies' \
//    -H 'Authorization: Bearer eyJhbGciOiJIUzI1Ni....heEAcZssEZOHkAhx-sbeUrHn8_YDcHY' \
//    --data-raw '{ \
//      "email": "newmovie@sabio.com", \
//      "password": "pass", \
//      "first_name": "firstname", \
//      "last_name": "lastname" \
//    }' \
// Response json data "createdMovieResponse"
func (c *movieController) CreateMovie(ctx echo.Context) (err error) {

	return ctx.JSON(http.StatusCreated, response{})
}
