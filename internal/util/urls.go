package util

import (
	"net/url"

	"github.com/pkg/errors"
)

func NormaliseURL(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if u.Scheme == "https" {
		u.Scheme = "http"
	}
	return u.String(), nil
}
