package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"time"

	"github.com/kevinburke/handlers"
	"github.com/kevinburke/rest"
	"github.com/saintpete/logrole/assets"
)

var errorTemplate *template.Template

func init() {
	base := string(assets.MustAsset("templates/base.html"))
	templates := template.Must(template.New("base").Option("missingkey=error").Funcs(funcMap).Parse(base))

	terror := template.Must(templates.Clone())
	errorTpl := string(assets.MustAsset("templates/errors.html"))
	errorTemplate = template.Must(terror.Parse(errorTpl))
}

type errorData struct {
	baseData
	Title       string
	Description string
	Mailto      *mail.Address
}

type errorServer struct {
	Mailto *mail.Address
}

func (e *errorServer) Serve401(w http.ResponseWriter, r *http.Request) {
	data := &errorData{
		Title:       "Unauthorized",
		Description: "Please enter your credentials to access this page.",
		Mailto:      e.Mailto,
	}
	domain := rest.CtxDomain(r)
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, domain))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(401)
	data.Start = time.Now()
	if err := render(w, errorTemplate, "base", data); err != nil {
		handlers.Logger.Error("Error rendering error template", "err", err)
	}
}

func (e *errorServer) Serve403(w http.ResponseWriter, r *http.Request) {
	data := &errorData{
		Title:       "Forbidden",
		Description: "You don't have permission to access this page",
		Mailto:      e.Mailto,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(403)
	data.Start = time.Now()
	if err := render(w, errorTemplate, "base", data); err != nil {
		handlers.Logger.Error("Error rendering error template", "err", err)
	}
}

func (e *errorServer) Serve404(w http.ResponseWriter, r *http.Request) {
	data := &errorData{
		Title:       "Page Not Found",
		Description: "We couldn't find what you were looking for. Sorry about that.",
		Mailto:      e.Mailto,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)
	data.Start = time.Now()
	if err := render(w, errorTemplate, "base", data); err != nil {
		handlers.Logger.Info("Error rendering error template", "err", err)
	}
}

func (e *errorServer) Serve405(w http.ResponseWriter, r *http.Request) {
	data := &errorData{
		Title:       "Method not allowed",
		Description: fmt.Sprintf("You can't make a %s request to this page.", r.Method),
		Mailto:      e.Mailto,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(405)
	data.Start = time.Now()
	if err := render(w, errorTemplate, "base", data); err != nil {
		handlers.Logger.Info("Error rendering error template", "err", err)
	}
}

func (e *errorServer) Serve500(w http.ResponseWriter, r *http.Request) {
	data := &errorData{
		Title:       "Server Error",
		Description: "Unexpected server error when serving your request. Please refresh the page and try again.",
		Mailto:      e.Mailto,
	}
	err := rest.CtxErr(r)
	handlers.Logger.Error("Server error", "code", 500, "method", r.Method, "path", r.URL.Path, "err", err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(500)
	data.Start = time.Now()
	if err := render(w, errorTemplate, "base", data); err != nil {
		handlers.Logger.Error("Error rendering error template", "err", err)
	}
}

func registerErrorHandlers(e *errorServer) {
	rest.RegisterHandler(401, e.Serve401)
	rest.RegisterHandler(403, e.Serve403)
	rest.RegisterHandler(404, e.Serve404)
	rest.RegisterHandler(405, e.Serve405)
	rest.RegisterHandler(500, e.Serve500)
}
