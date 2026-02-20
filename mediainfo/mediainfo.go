package mediainfo

import (
	"encoding/json"
	"log/slog"
	"os/exec"
)

type StreamKind int

const (
	StreamGeneral StreamKind = iota
	StreamVideo
	StreamAudio
	StreamText
	StreamOther
	StreamImage
	StreamMenu
	streamMax
)

type InfoKind int

const (
	// Unique name of parameter
	InfoName InfoKind = iota
	// Value of parameter
	InfoText
	// Unique name of measure unit of parameter
	InfoMeasure
	// See infooptions_t
	InfoOptions
	// Translated name of parameter
	InfoNameText
	// Translated name of measure unit
	InfoMeasureText
	// More information about the parameter
	InfoInfo
	// Information : how data is found
	InfoHowTo
	infoMax
)

func Handle(filePath string) ([]map[string]any, error) {
	slog.Debug("opening file", "filename", filePath)

	cmd := exec.Command("mediainfo", "--Output=JSON", "-f", filePath)
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	m := struct {
		Media struct {
			Track []map[string]any `json:"track"`
		} `json:"media"`
	}{}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &m)
	if err != nil {
		return nil, err
	}

	return m.Media.Track, nil
}
