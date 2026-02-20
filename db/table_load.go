package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func (t *Table[V]) load(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		bytes := scanner.Bytes()

		slog.Debug("table load scan", "v", fmt.Sprintf("%T", *new(V)), "line", string(bytes))

		var v V

		err = json.Unmarshal(bytes, &v)
		if err != nil {
			slog.Warn("issue unmarshalling", "error", err)

			continue
		}

		if t.exists(v) {
			slog.Debug("table load unique", "v", fmt.Sprintf("%T", *new(V)), "line", string(bytes))

			continue
		}

		if t.append(v) == nil {
			slog.Warn("table load append with nil return", "v", fmt.Sprintf("%T", *new(V)), "line", string(bytes))
		}
	}

	return nil
}

func (t *Table[V]) Load(paths ...string) error {
	for _, path := range paths {
		slog.Debug("table load", "v", fmt.Sprintf("%T", *new(V)), "path", path)

		err := t.load(path)
		if err != nil {
			return err
		}
	}

	return nil
}
