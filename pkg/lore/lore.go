package lore

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

func Parse(url string) (string, error) {
	err := errors.New("lore: not a valid url")

	if !strings.HasPrefix(url, "https://lore.kernel.org/") {
		return "", err
	}

	if strings.HasSuffix(url, "/t.mbox.gz") {
		return url, nil
	}

	pos := strings.LastIndex(url, "/T")
	if pos < 0 {
		return url, err
	}

	url = url[:pos] + "/t.mbox.gz"
	return url, nil
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
