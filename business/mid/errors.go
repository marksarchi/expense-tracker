package mid

import (
	"log"
	"net/http"

	"github.com/sarchimark/expense-tracker/foundation/web"
)

func Errors(log *log.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			if err := handler(w, r); err != nil {
				log.Printf(" ERROR: %v", err)
				if err := web.RespondError(w, err); err != nil {
					return err

				}

			}
			return nil

		}
		return h
	}
	return m
}
