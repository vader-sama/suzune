package util

import (
	"net/url"
	"path"
	"path/filepath"
)

func ResolveOutputPath(urlStr, outputFlag string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(outputFlag)
	if ext != "" {
		return outputFlag, nil
	}

	filename := path.Base(u.Path)
	return filepath.Join(outputFlag, filename), nil
}
