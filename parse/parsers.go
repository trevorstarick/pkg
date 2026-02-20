package parse

import (
	"errors"
	"strconv"
	"strings"
)

const sixteenbynine = 16.0 / 9.0

//nolint:mnd // magic numbers are fine here
func Resolution(s string) (float64, error) {
	switch strings.ToLower(s) {
	case "sd", "480p":
		return 480 * 480 * sixteenbynine, nil
	case "hd", "720p":
		return 720 * 720 * sixteenbynine, nil
	case "full hd", "fullhd", "fhd", "1080p":
		return 1080 * 1080 * sixteenbynine, nil
	case "quad hd", "quadhd", "qhd", "1440p", "2k":
		return 1440 * 1440 * sixteenbynine, nil
	case "ultra hd", "ultrahd", "uhd", "2160p", "4k":
		return 2160 * 2160 * sixteenbynine, nil
	case "5k":
		return 2880 * 2880 * sixteenbynine, nil
	case "8k":
		return 4320 * 4320 * sixteenbynine, nil
	}

	if before, ok := strings.CutSuffix(strings.ToLower(s), "mp"); ok {
		s = before

		megapixels, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, errors.New("strconv.ParseFloat")
		}

		megapixels *= 1_000_000

		return megapixels, nil
	}

	if strings.Contains(s, "x") {
		parts := strings.SplitN(s, "x", 2)
		widthStr, heightStr := parts[0], parts[1]

		width, err := strconv.ParseFloat(widthStr, 64)
		if err != nil {
			return 0, errors.New("strconv.ParseFloat")
		}

		height, err := strconv.ParseFloat(heightStr, 64)
		if err != nil {
			return 0, errors.New("strconv.ParseFloat")
		}

		return width * height, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("strconv.ParseFloat")
	}

	return f, nil
}

func Bits(s, suffix string) (int, error) {
	s = strings.TrimSuffix(strings.ToLower(s), suffix)

	bits, err := Figs(s)
	if err != nil {
		return 0, errors.New("parseBits")
	}

	return int(bits), nil
}

func Figs(s string) (float64, error) {
	mul := 1.0

	s = strings.ToLower(s)

	switch {
	case strings.HasSuffix(s, "k"):
		s = strings.TrimSuffix(s, "k")
		mul = 1_000
	case strings.HasSuffix(s, "m"):
		s = strings.TrimSuffix(s, "m")
		mul = 1_000_000
	case strings.HasSuffix(s, "g"):
		s = strings.TrimSuffix(s, "g")
		mul = 1_000_000_000
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("strconv.Atoi")
	}

	return mul * f, nil
}

func Megapixels(s string) (float64, error) {
	return Figs(strings.TrimSuffix(s, "p"))
}

func Filesize(s string) (int, error) {
	return Bits(s, "b")
}

func Bitrate(s string) (int, error) {
	return Bits(s, "bps")
}
