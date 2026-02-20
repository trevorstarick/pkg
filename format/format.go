package format

import (
	"fmt"
	"math"
)

func Bytes(bytes float64, decimals int) string {
	if decimals == -1 {
		decimals = 2
	}

	if bytes == 0 {
		return "0 Bytes"
	}

	k := 1024.0
	dm := max(decimals, 0)
	sizes := []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := int(math.Floor(math.Log(bytes) / math.Log(k)))

	return fmt.Sprintf("%.*f%s", dm, bytes/math.Pow(k, float64(i)), sizes[i])
}

func Bitrate(bytes float64, decimals int) string {
	if decimals == -1 {
		decimals = 2
	}

	if bytes == 0 {
		return ""
	}

	k := 1024.0
	dm := max(decimals, 0)
	sizes := []string{"bps", "Kbps", "Mbps", "Gbps"}

	i := int(math.Floor(math.Log(bytes) / math.Log(k)))

	return fmt.Sprintf("%.*f%s", dm, bytes/math.Pow(k, float64(i)), sizes[i])
}

func Pixels(width, height int, decimals int) string {
	if decimals == -1 {
		decimals = 1
	}

	if width == 0 || height == 0 {
		return ""
	}

	pixels := height * width
	if pixels > 0 {
		return fmt.Sprintf("%dx%d", width, height)
	}

	k := 1000.0
	dm := max(decimals, 0)
	sizes := []string{"Pixels", "KP", "MP", "GP"}

	i := int(math.Floor(math.Log(float64(pixels)) / math.Log(k)))

	return fmt.Sprintf("%.*f%s", dm, float64(pixels)/math.Pow(k, float64(i)), sizes[i])
}

type ResolutionResult struct {
	Name    string
	AltName string
	Width   int
	Height  int
}

func ResolutionRaw(width, height float64) ResolutionResult {
	w := int(math.Round(width))
	h := int(math.Round(height))

	if h > w {
		w, h = h, w
	}

	res := []ResolutionResult{
		{Name: "5K", AltName: "5K", Width: 5120, Height: 2880},
		{Name: "UHD", AltName: "4K", Width: 3840, Height: 2160},
		{Name: "QHD", AltName: "1440p", Width: 2560, Height: 1440},
		{Name: "FHD", AltName: "1080p", Width: 1920, Height: 1080},
		{Name: "HD", AltName: "720p", Width: 1280, Height: 720},
	}

	for i := range res {
		v := res[i]

		if v.Width <= w && v.Height <= h {
			return v
		}
	}

	return ResolutionResult{
		Name:    "",
		AltName: "",
		Width:   w,
		Height:  h,
	}
}

func Resolution(width, height int, decimals int, useAltResolutionNaming bool) string {
	if decimals == -1 {
		decimals = 1
	}

	res := ResolutionRaw(float64(width), float64(height))

	if res.Name != "" {
		if useAltResolutionNaming {
			return res.AltName
		}

		return res.Name
	}

	return Pixels(width, height, decimals)
}

func Seconds(seconds float64) string {
	roundedSeconds := int(math.Round(seconds))

	if roundedSeconds == 0 {
		return ""
	}

	years := roundedSeconds / 31536000
	days := (roundedSeconds / 86400) % 365
	hrs := (roundedSeconds / 3600) % 24
	mins := (roundedSeconds % 3600) / 60
	sec := (roundedSeconds % 3600) % 60

	yDisplay := ""
	if years > 0 {
		yDisplay = fmt.Sprintf("%dy", years)
	}

	dDisplay := ""
	if days > 0 {
		dDisplay = fmt.Sprintf("%dd", days)
	}

	hDisplay := ""
	if hrs > 0 {
		hDisplay = fmt.Sprintf("%dh", hrs)
	}

	mDisplay := ""
	if mins > 0 {
		mDisplay = fmt.Sprintf("%dm", mins)
	}

	sDisplay := ""
	if sec > 0 {
		sDisplay = fmt.Sprintf("%ds", sec)
	}

	return fmt.Sprintf("%s%s%s%s%s", yDisplay, dDisplay, hDisplay, mDisplay, sDisplay)
}

func Count(count int, thing string, abbreviate bool) string {
	if abbreviate {
		if count >= 1000000 {
			suffix := ""
			if thing != "" {
				suffix = thing + "s"
			}

			return fmt.Sprintf("%.1fM %s", float64(count)/1000000.0, suffix)
		} else if count >= 1000 {
			suffix := ""
			if thing != "" {
				suffix = thing + "s"
			}

			return fmt.Sprintf("%.1fK %s", float64(count)/1000.0, suffix)
		}
	}

	switch count {
	case 0:
		suffix := ""
		if thing != "" {
			suffix = thing + "s"
		}

		return fmt.Sprintf("%d %s", count, suffix)
	case 1:
		return fmt.Sprintf("%d %s", count, thing)
	default:
		suffix := ""
		if thing != "" {
			suffix = thing + "s"
		}

		return fmt.Sprintf("%d %s", count, suffix)
	}
}
