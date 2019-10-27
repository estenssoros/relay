package goflow

import (
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/estenssoros/goflow/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/leaanthony/mewn"
	"github.com/leaanthony/mewn/lib"
	"github.com/pkg/errors"
)

// Webserver handles the webserver
type Webserver struct {
	Dags map[string]*DAG
}

// NewWebserver creates a new webserver from a dag map
func NewWebserver(dags map[string]*DAG) *Webserver {
	return &Webserver{dags}
}

// Routes applies routes to echo
func (w *Webserver) Routes(e *echo.Echo) {
	e.GET("/", func(e echo.Context) error {
		return e.Render(http.StatusOK, "index.html", nil)
	})
}

// Serve serves the web endopoints
func (w *Webserver) Serve() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	templateGroup := mewn.Group("./templates")
	t := &Template{
		templates: templateGroup,
	}
	e.Renderer = t
	w.Routes(e)

	e.Logger.Fatal(
		e.Start(":" + strconv.FormatInt(config.DefaultConfig.Webserver.Port, 10)),
	)

	return nil
}

// Template holdes templates
type Template struct {
	templates *lib.FileGroup
}

// Render implements the echo render function
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	ts := t.templates.String(name)
	if ts == "" {
		return errors.Errorf("missing template: %s", name)
	}
	temp := template.Must(template.New(name).Parse(ts))
	return temp.ExecuteTemplate(w, name, data)
}
