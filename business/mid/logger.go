package mid

import (
	"log"
	"net/http"
	"time"

	"github.com/sarchimark/expense-tracker/foundation/web"
)

func Logger(log *log.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {

			t := time.Now()
			log.Printf("started: %s %s -> %s", r.Method, r.URL.Path, r.RemoteAddr)

			//Call the next handler.
			err := handler(w, r)

			log.Printf("completed: %s %s  -> %s  (%s)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(t))

			return err

		}
		return f

	}
	return m

}
