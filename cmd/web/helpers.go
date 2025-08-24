package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// serverError helper writes log entry at Error level and sends generic 500 internal server error.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method        = r.Method
		uri           = r.URL.RequestURI()
		trace         = debug.Stack()
		extraLogsInfo = []any{
			slog.String("method", method),
			slog.String("uri", uri),
			slog.String("stack", string(trace)),
		}
	)

	app.logger.Error(err.Error(), extraLogsInfo...)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError helper sends a specific status code and corresponding description.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
