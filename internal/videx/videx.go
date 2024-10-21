package videx

import (
	"fmt"
)

// ExtractVideos checks the URL and calls appropriate extractor funcs
func ExtractVideos(url string) ([]Video, error) {
	youtubeExtractor := GetExtractor(url)
	if youtubeExtractor == nil {
		return nil, fmt.Errorf("[ERROR] no extractor available for URL: %s", url)
	}
	return youtubeExtractor.ExtractVideos(url)
}

// ExtractPlayback checks the URL and calls the appropriate extractor function.
func ExtractPlayback(url string) (string, error) {
	extractor := GetExtractor(url)
	if extractor == nil {
		return "", fmt.Errorf("[ERROR] no extractor available for URL: %s", url)
	}
	return extractor.ExtractPlayback(url)
}
