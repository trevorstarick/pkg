package htmx

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func (htmx *HTMX) partialLoader(name string) (template.HTML, error) {
	if name == "" {
		slog.Debug("load function: empty name")

		return template.HTML(""), nil
	}

	var filename string
	for _, dir := range append(htmx.TemplateDirs, "views") {
		filename = filepath.Join(dir, "components", name+".html")

		if _, err := os.Stat(filename); err == nil {
			break
		}
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return template.HTML(""), fmt.Errorf("read partial file %s: %w", filename, err)
	}

	//nolint:gosec // bytes are read from a file, so they should be safe to use as HTML
	return template.HTML(strings.TrimSuffix(string(bytes), "\n")), nil
}

// convertRequestToTemplatePaths generates a list of potential template paths based on the request URL and pattern. It prioritizes more specific paths first, followed by more general ones, and finally includes a default index path.
// For example, if the request URL is "/foo/bar" and the pattern is "/foo/*", it will generate paths like "ssr/foo/bar.html", "ssr/foo/index.html", "ssr/foo.html", and "ssr/index.html". The function also checks for the "Hx-Request" header to determine whether to use the "ssr/" or "pages/" directory for the templates.
// The generated paths are sorted in descending order of specificity (i.e., paths with more segments are prioritized over those with fewer segments) to ensure that the most specific template is attempted first when rendering the response.
func (htmx *HTMX) convertRequestToTemplatePaths(r *http.Request) []string {
	var part string

	if r.Header.Get("Hx-Request") == "true" {
		part = "ssr/"
	} else {
		part = "pages"
	}

	path := r.URL.Path
	paths := []string{}
	for {
		p := filepath.Join(part, strings.TrimSuffix(path[1:], "/"))
		paths = append(paths, p+".html")
		if filepath.Dir(p) != "/" {
			paths = append(paths, filepath.Join(filepath.Dir(p), "index.html"))
		} else {
			break
		}

		path = filepath.Dir(path)

		if path == "/" {
			break
		}
	}

	paths = append(paths, filepath.Join(part, strings.TrimSuffix(r.Pattern[1:], "/")+".html"))

	slices.SortFunc(paths, func(i, j string) int {
		return strings.Count(j, "/") - strings.Count(i, "/")
	})

	paths = append(paths, filepath.Join(part, "index.html"))

	uniquePathsMap := map[string]struct{}{}
	uniquePaths := []string{}
	for _, p := range paths {
		if _, exists := uniquePathsMap[p]; !exists {
			uniquePathsMap[p] = struct{}{}

			for _, location := range append(htmx.TemplateDirs, []string{
				"views",
			}...) {
				uniquePaths = append(uniquePaths, filepath.Join(location, p))
			}
		}
	}

	return uniquePaths
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

func (htmx *HTMX) loadPartials(data string) (string, error) {
	re := regexp.MustCompile(`{{ load "(.*)" }}`)

	for re.MatchString(data) {
		loads := re.FindAllStringSubmatch(data, -1)

		for _, load := range loads {
			partialName := load[1]

			partialContent, err := htmx.partialLoader(partialName)
			if err != nil {
				return "", fmt.Errorf("load partial %s: %w", partialName, err)
			}

			data = strings.ReplaceAll(data, load[0], string(partialContent))
		}
	}

	return data, nil
}

func (htmx HTMX) getTemplateBody(path string) (*template.HTML, error) {
	data, err := readTemplateFile(path)
	if err != nil {
		return nil, err
	}

	data, err = htmx.loadPartials(data)
	if err != nil {
		return nil, fmt.Errorf("load partials: %w", err)
	}

	//nolint:gosec // data is read from a file, so it should be safe to use as HTML
	templateContent := template.HTML(data)

	return &templateContent, nil
}
