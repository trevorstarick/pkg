package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func (t *Table[V]) Save(path string) error {
	var f io.WriteCloser

	if !t.changedSinceLastSave.Load() {
		return nil
	}

	filename := filepath.Base(path) + time.Now().Format("_20060102_150405")
	dir := filepath.Dir(path)

	if os.Getenv("DEBUG") == "true" {
		f = os.Stdout
	} else {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}

		slog.Debug("table save", "v", fmt.Sprintf("%T", *new(V)), "path", path, "filename", filename)

		f, err = os.OpenFile(filepath.Join(dir, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
	}

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	slog.Info("table save", "v", fmt.Sprintf("%T", *new(V)), "len", t.data.Len())
	t.data.SortedRange(func(_ any, p *V) bool {
		if p == nil {
			slog.Error("nil value encountered during save", "v", fmt.Sprintf("%T", *new(V)))

			return true
		}

		v := *p

		b, err := json.Marshal(v)
		if err != nil {
			slog.Error("table save marshal", "v", fmt.Sprintf("%T", v), "error", err)

			return true
		}

		_, err = writer.Write(b)
		if err != nil {
			slog.Error("table save write", "v", fmt.Sprintf("%T", v), "error", err)

			return true
		}

		_, err = writer.WriteString("\n")
		if err != nil {
			slog.Error("table save write newline", "v", fmt.Sprintf("%T", v), "error", err)

			return true
		}

		return true
	})

	err := writer.Flush()
	if err != nil {
		slog.Error("table save flush", "error", err)
	}

	if f != os.Stdout {
		err := f.Close()
		if err != nil {
			slog.Error("table save close", "error", err)
		}

		err = os.Rename(filepath.Join(dir, filename), path)
		if err != nil {
			slog.Error("table save rename", "error", err)
		}
	}

	t.changedSinceLastSave.Store(false)

	return nil
}
