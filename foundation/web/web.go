package web

import (
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers
type App struct {
	mux      *chi.Mux
	mw       []Middleware
	shutdown chan os.Signal
}

//A Handler is  funcrion type that handles a http request
type Handler func(w http.ResponseWriter, r *http.Request) error

//NewApp creates an App instance that handle a set of routes
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	mux := chi.NewRouter()

	return &App{
		mux:      mux,
		shutdown: shutdown,
		mw:       mw,
	}
}

//SignalShutdown is used to gracefully shutdown the app when an issue is identified
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

//ServeHTTP implements the http.Handler interface.Its the entry point of all http traffic
//it passes the request to the mux which then routes to the appropriate Handler
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

//Handle sets a handler function to the application server mux
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	handler = wrapMiddleware(mw, handler)

	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		if err := handler(w, r); err != nil {
			a.SignalShutdown()
			return

		}
	}

	a.mux.MethodFunc(method, path, h)

}
