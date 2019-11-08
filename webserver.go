package relay

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/estenssoros/relay/config"
	"github.com/estenssoros/relay/db"
	"github.com/estenssoros/relay/models"
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
func (w *Webserver) Routes(c *echo.Echo) {
	group := c.Group("/api")
	group.GET("/dags", func(c echo.Context) error {
		conn := db.Connection
		dags := []*models.DAG{}
		if err := conn.Find(&dags).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, dags)
	})
	group.POST("/kill", func(c echo.Context) error {
		ctx := c.Get("appContext").(context.Context)
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		fmt.Println("cancelled!")
		return c.JSON(http.StatusOK, "killed")
	})
	group.POST("/dag-toggle", func(c echo.Context) error {
		req := &struct {
			DagID  string `json:"dagID"`
			Paused bool   `json:"paused"`
		}{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		dag := &models.DAG{ID: req.DagID}
		if err := db.Connection.Model(dag).Update("IsPaused", req.Paused).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return nil
	})
}

type mewnFileServer struct {
	group *lib.FileGroup
}

func (fs *mewnFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
}

// Serve serves the web endopoints
func (w *Webserver) Serve(ctx context.Context) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// templateGroup := mewn.Group("./templates")
	// t := &Template{
	// 	templates: templateGroup,
	// }
	// e.Renderer = t
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("appContext", ctx)
			return next(c)
		}
	})

	w.Routes(e)

	templateGroup := mewn.Group("./templates")
	e.GET("/", func(e echo.Context) error {
		fileStr := templateGroup.String("index.html")
		return e.Blob(http.StatusOK, "text/html", []byte(fileStr))
	})

	e.GET("/favicon.ico", func(e echo.Context) error {
		fileStr := templateGroup.String("favicon.ico")
		return e.Blob(http.StatusOK, "image/ico", []byte(fileStr))
	})

	e.GET("/public/:fileName", func(e echo.Context) error {
		var fileContentType string

		fileStr := templateGroup.String(e.Param("fileName"))

		switch filepath.Ext(e.Param("fileName")) {
		case ".css":
			fileContentType = "text/css"
		default:
			reader := bytes.NewReader([]byte(fileStr))
			fileHeader := make([]byte, 512)
			reader.Read(fileHeader)
			fileContentType = http.DetectContentType(fileHeader)
		}
		return e.Blob(http.StatusOK, fileContentType, []byte(fileStr))
	})

	c := make(chan error)
	go func() {
		err := e.Start(":" + strconv.FormatInt(config.DefaultConfig.Webserver.Port, 10))
		if err != nil {
			c <- err
		}
		close(c)
	}()
	for {
		select {
		case err := <-c:
			log.Fatal(err)
		case <-ctx.Done():
			return nil
		}
	}
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
