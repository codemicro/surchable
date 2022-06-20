package util

import (
	"github.com/PuerkitoBio/purell"
	"github.com/pkg/errors"
)

func NormaliseURL(inputURL string) (string, error) {
	normalisedURL, err := purell.NormalizeURLString(
		inputURL,
		purell.FlagsUsuallySafeGreedy|purell.FlagForceHTTP,
	)
	return normalisedURL, errors.WithStack(err)
}
