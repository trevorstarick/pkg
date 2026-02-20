package funcs

import (
	"encoding/json"
	"log/slog"
)

func JSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		slog.Error("json marshal", "error", err)

		return ""
	}

	return string(b)
}
