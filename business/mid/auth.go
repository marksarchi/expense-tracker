package mid

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			//Expecting bearer token
			authStr := r.Header.Get("authorization")

			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)

			}

			//Validate the token
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			//claims := 2

			//Add claims to the context so they can be retrieved later
			ctx := context.WithValue(r.Context(), auth.Key, claims)

			return handler(w, r.WithContext(ctx))

		}
		return f
	}
	return m
}
