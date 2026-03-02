package htmx

import (
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/trevorstarick/pkg/htmx/funcs"
)

// this is a global variable that holds the default template funcs, and is not
// modified after initialization, so it's safe to use as a global variable
//
//nolint:gochecknoglobals
var templateFuncs = map[string]any{
	"json": funcs.JSON,

	"spew": spew.Sdump,

	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"contains":   strings.Contains,
	"toLower":    strings.ToLower,
	"toUpper":    strings.ToUpper,
	"trimSpace":  strings.TrimSpace,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,

	"replace": strings.ReplaceAll,
}

// this is a global variable that holds the extra template funcs, and is
// modified by the SetTemplateFunc function, so it's safe to use as a global variable
//
//nolint:gochecknoglobals
var extraTemplateFuncs = map[string]any{}

func SetTemplateFunc(name string, fn any) {
	extraTemplateFuncs[name] = fn
}

func getTemplateFuncs(r *http.Request) map[string]any {
	funcs := extraTemplateFuncs
	// call maps.Copy with templateFuncs so we guarantee extraTemplateFuncs won't
	//   override any of the default funcs, but still allow extraTemplateFuncs to add new funcs
	maps.Copy(funcs, templateFuncs)

	if r != nil {
		funcs["ctx"] = r.Context().Value
		funcs["queryParams"] = func() url.Values { return r.URL.Query() }
		funcs["isBoosted"] = func() bool { return r.Header.Get("Hx-Boosted") == "true" }
	} else {
		fn := func() { slog.Warn("func not supported when request is nil") }

		funcs["ctx"] = fn
		funcs["queryParams"] = fn
		funcs["isBoosted"] = fn
	}

	return funcs
}
