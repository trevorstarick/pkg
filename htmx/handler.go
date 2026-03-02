package htmx

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func (htmx *HTMX) TemplateHandlerFunc(data ...any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := htmx.TemplateHandler(data...)(w, r)
		if err != nil {
			slog.Error("template handler", "error", err)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (htmx *HTMX) TemplateHandler(data ...any) func(w http.ResponseWriter, r *http.Request) error {
	if len(data) == 0 {
		data = append(data, nil)
	} else if len(data) > 1 {
		d := data

		data = []any{d}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		funcs := getTemplateFuncs(r)
		pathNames := htmx.convertRequestToTemplatePaths(r)

		templateBody, err := htmx.getTemplateBody(path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			templateData, err := os.ReadFile("views/404.html")
			if err != nil {
				return err
			}

			tmpl, err := template.New("").Funcs(funcs).Parse(string(templateData))
			if err != nil {
				return err
			}

			w.WriteHeader(http.StatusNotFound)

			err = tmpl.Execute(w, map[string]any{
				"title": "Page Not Found",
				"body":  template.HTML("<h1>404 - Page Not Found</h1>"),
			})
			if err != nil {
				return err
			}

			return nil
		}

		var t string

		if r.Header.Get("Hx-Request") == "true" || r.Header.Get("Hx-Boosted") == "true" {
			t = string(*templateBody)
		} else {
			slog.Debug("render template", "path", r.URL.Path, "pattern", r.Pattern)

			templateData, err := os.ReadFile("views/200.html")
			if err != nil {
				return err
			}

			t = string(templateData)
			t = strings.ReplaceAll(t, "{{ body }}", string(*templateBody))
		}

		w.Header().Set("Content-Type", "text/html")

		tmpl, err := template.New("").Funcs(funcs).Parse(t)
		if err != nil {
			return err
		}

		err = tmpl.Execute(w, data[0])
		if err != nil {
			return err
		}

		return nil
	}
}
