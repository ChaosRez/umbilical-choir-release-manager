package repository

import (
	"os"
)

const releaseFile = "releases/release.yml"

func GetLatestRelease() (string, error) {
	if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
		return "", err
	}
	return releaseFile, nil
}
func NewReleaseExists() bool {
	_, err := os.Stat(releaseFile)
	return !os.IsNotExist(err)
}
