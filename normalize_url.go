package main

import (
	"errors"
	"net/url"
	"strings"

)


func normalizeURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("invalid URL")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	path := strings.TrimSuffix(parsedURL.Path, "/")

	return parsedURL.Host + path, nil
}
