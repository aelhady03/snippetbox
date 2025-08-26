package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

// render renders a template with associated data.
// render the template in two stages: write them into buffer first
// if everything is good, write from the buffer to the response writer
// if something wrong, return a server error to the user.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// write it to the buffer first?
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// all is good? write then to the response writer
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		var invalidDecodeErr *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeErr) {
			panic(err)
		}
		return err
	}

	return nil
}
