package videx

import "strings"

// Extractor defines the interface for video extractors
type Extractor interface {
	ExtractVideos(url string) ([]Video, error)
	ExtractPlayback(url string) (string, error)
}

// YouTubeExtractor implements the Extractor interface
type YouTubeExtractor struct{}

func newYouTubeExtractor() *YouTubeExtractor {
	return &YouTubeExtractor{}
}


func (yt *YouTubeExtractor) ExtractVideos(url string) ([]Video, error) {
	src, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	var index int
	var videos []Video

	for {
		var compact, channel bool

		tempIndex := strings.Index(src[index:], "compactVideoRenderer\":{\"videoId")
		if tempIndex == -1 {
			tempIndex = strings.Index(src[index:], "videoRenderer\":{\"videoId")
			if tempIndex == -1 {
				tempIndex = strings.Index(src[index:], "VideoRenderer\":{\"videoId")
				if tempIndex == -1 {
					break // No more video renderers found, end loop
				}
				channel = true
			}
		} else {
			compact = true
		}

		index += tempIndex

		video := Video{}

		// Safely extract video URL
		videoIDStart := strings.Index(src[index:], "\"videoId\":\"")
		if videoIDStart == -1 || index+videoIDStart+len("\"videoId\":\"")+11 > len(src) {
			break // Invalid or missing videoId, break the loop
		}
		index += videoIDStart + len("\"videoId\":\"")
		video.URL = "/watch?v=" + src[index:index+11]

		// Safely extract thumbnail
		thumbnailStart := strings.Index(src[index:], "\"url\":\"")
		if thumbnailStart == -1 {
			break // Missing thumbnail, break loop
		}
		index += thumbnailStart + len("\"url\":\"")
		thumbnailEnd := strings.Index(src[index:], "\"")
		if thumbnailEnd == -1 || index+thumbnailEnd > len(src) {
			break // Invalid thumbnail, break loop
		}
		video.Thumbnail = src[index:index+thumbnailEnd]

		// Safely extract video title
		var titleStart int
		if compact {
			titleStart = strings.Index(src[index:], "simpleText\":\"")
			if titleStart == -1 {
				break // No title found
			}
			index += titleStart + len("simpleText\":\"")
		} else {
			titleStart = strings.Index(src[index:], "\"title\":{\"runs\":[{\"text\":\"")
			if titleStart == -1 {
				break // No title found
			}
			index += titleStart + len("\"title\":{\"runs\":[{\"text\":\"")
		}
		titleEnd := strings.Index(src[index:], "\"")
		if titleEnd == -1 || index+titleEnd > len(src) {
			break // Invalid title, break loop
		}
		video.Title = src[index : index+titleEnd]

		// Extract channel info if not a channel itself
		if !channel {
			bylineStart := strings.Index(src[index:], "longBylineText")
			if bylineStart != -1 {
				index += bylineStart + len("longBylineText\":{\"runs\":[{\"text\":\"")
				bylineEnd := strings.Index(src[index:], "\"")
				if bylineEnd != -1 && index+bylineEnd <= len(src) {
					video.Channel = src[index : index+bylineEnd]
				}
			}

			lengthTextStart := strings.Index(src[index:], "lengthText")
			if lengthTextStart != -1 {
				index += lengthTextStart
				simpleTextStart := strings.Index(src[index:], "simpleText\":\"")
				if simpleTextStart != -1 {
					index += simpleTextStart + len("simpleText\":\"")
					lengthEnd := strings.Index(src[index:], "\"")
					if lengthEnd != -1 && index+lengthEnd <= len(src) {
						video.Length = src[index : index+lengthEnd]
					}
				}
			}
		}

		// Append extracted video to the list
		videos = append(videos, video)
	}

	return videos, nil
}


func (yt *YouTubeExtractor) ExtractPlayback(url string) (string, error) {
	src, err := makeRequest(url)
	if err != nil {
		return "", err
	}

	index := strings.Index(src, "googlevideos.com\\/videosplayback")
	if index != -1 {
		index -= 30
		index = strings.Index(src[index:], "https:") + index

		watchURL := src[index : strings.Index(src[index:], "\"")+index]

		// replacing `\\/` with `/`
		for strings.Contains(watchURL, "\\/") {
			watchURL = strings.Replace(watchURL, "\\/", "/", -1)
		}

		// replacing `\\\\u0026` with `&`
		for strings.Contains(watchURL, "\\\\u0026") {
			watchURL = strings.Replace(watchURL, "\\\\u0026", "&", -1)
		}
		return watchURL, nil
	}
	return "", nil
}

// GetExtractor is a factory func that returns the appropriate Extractor
// interface based on the URL
func GetExtractor(url string) Extractor {
	if strings.Contains(url, "youtube.com") {
		return newYouTubeExtractor()
	}
	// additional extractors for other platforms go here
	return nil
}
