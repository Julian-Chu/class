package handlers

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

func health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
