package videx

import (
	"regexp"
	"strconv"
	"strings"
)

type Video struct {
	URL        string
	Title      string
	Length     string
	Thumbnail  string
	Channel    string
	ChannelURL string
}

// CleanupTitle cleans up the video title by replacing HTML entities
// and unescaping unicode characters.
func (v *Video) CleanupTitle() {
	replacements := map[string]string{
		"&quot;": "\"",
		"&amp;":  "&",
		"–":      "-",
		"…":      "...",
		"’":      "'",
	}

	for entity, char := range replacements {
		v.Title = strings.ReplaceAll(v.Title, entity, char)
	}

	// unescape unicode chars encoded as \uXXX
	re := regexp.MustCompile(`\\u([0-9a-fA-F]){4}`)
	v.Title = re.ReplaceAllStringFunc(v.Title, func(match string) string {
		code := match[2:]
		r, err := runeFromHex(code)
		if err != nil {
			return match
		}
		return string(r)
	})
}

// runeFromHex converts a hex str to a rune
func runeFromHex(hexStr string) (rune, error) {
	code, err := strconv.ParseInt(hexStr, 16, 32)
	if err != nil {
		return 0, err
	}
	return rune(code), nil
}
