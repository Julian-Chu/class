package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ardanlabs/service/foundation/web"
)

// Logger ...
func Logger(log *log.Logger) web.Middleware {
	m := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return nil // web.NewShutdownError("web value missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				v.TraceID, v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Now),
			)

			return err
		}

		return h
	}

	return m
}
