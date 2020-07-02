// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/service/business/auth"
	"github.com/ardanlabs/service/business/mid"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Error(log), mid.Metrics(), mid.Panic(log))

	c := check{
		build: build,
		db:    db,
	}
	app.Handle(http.MethodGet, "/health", c.health)

	p := productHandlers{
		db: db,
	}
	app.Handle(http.MethodGet, "/products", p.list, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/products/:id", p.retrieve, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/products", p.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/products/:id", p.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/products/:id", p.delete, mid.Authenticate(a))

	return app
}
