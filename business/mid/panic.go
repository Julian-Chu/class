package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ardanlabs/service/foundation/web"
	"github.com/pkg/errors"
)

// Panic ...
func Panic(log *log.Logger) web.Middleware {

	m := func(after web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)
					log.Printf("%s :\n%s", v.TraceID, debug.Stack())
				}
			}()

			return after(ctx, w, r)
		}

		return h
	}

	return m
}
