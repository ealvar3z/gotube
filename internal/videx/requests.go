package videx

import (
	"io"
	"net/http"
)

func makeRequest(url string) (string, error) {
	request, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer request.Body.Close()

	result, err := io.ReadAll(request.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
