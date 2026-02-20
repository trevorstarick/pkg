package htmx

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func partialLoader(name string) (template.HTML, error) {
	if name == "" {
		slog.Debug("load function: empty name")

		return template.HTML(""), nil
	}

	filename := filepath.Join("views/components/", name+".html")

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return template.HTML(""), fmt.Errorf("read partial file %s: %w", filename, err)
	}

	//nolint:gosec // bytes are read from a file, so they should be safe to use as HTML
	return template.HTML(strings.TrimSuffix(string(bytes), "\n")), nil
}

func convertRequestToTemplatePath(r *http.Request) string {
	var part string

	if r.Header.Get("Hx-Request") == "true" {
		part = "views/ssr/"
	} else {
		part = "views/pages"
	}

	pathName := "index"

	if r.Pattern != "" && r.Pattern != "/" {
		pathName = r.Pattern[1:]
	} else if r.URL.Path != "/" {
		pathName = r.URL.Path[1:]
	}

	pathName = strings.TrimSuffix(pathName, "/")

	return filepath.Join(part, pathName+".html")
}

func readTemplateFile(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		slog.Info("template not found", "path", path)

		return "", os.ErrNotExist
	}

	contentData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contentData), nil
}

func loadPartials(data string) (string, error) {
	re := regexp.MustCompile(`{{ load "(.*)" }}`)

	for re.MatchString(data) {
		loads := re.FindAllStringSubmatch(data, -1)

		for _, load := range loads {
			partialName := load[1]

			partialContent, err := partialLoader(partialName)
			if err != nil {
				return "", fmt.Errorf("load partial %s: %w", partialName, err)
			}

			data = strings.ReplaceAll(data, load[0], string(partialContent))
		}
	}

	return data, nil
}

func getTemplateBody(path string) (*template.HTML, error) {
	data, err := readTemplateFile(path)
	if err != nil {
		return nil, err
	}

	data, err = loadPartials(data)
	if err != nil {
		return nil, fmt.Errorf("load partials: %w", err)
	}

	//nolint:gosec // data is read from a file, so it should be safe to use as HTML
	templateContent := template.HTML(data)

	return &templateContent, nil
}
