package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

// Error ...
func Error(log *log.Logger) web.Middleware {

	m := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			if err := before(ctx, w, r); err != nil {

				log.Printf("%s : ERROR : %v", v.TraceID, err)

				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				if ok := web.IsShutdown(err); ok {
					// Validate if this is really something
					// we want to do.
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
