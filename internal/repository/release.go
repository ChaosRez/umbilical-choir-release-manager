package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

const latestReleaseFile = "releases/21/release.yml"

// TODO: check if both fns.zip file and release.yml file exist
func GetLatestRelease() (string, error) {
	if _, err := os.Stat(latestReleaseFile); os.IsNotExist(err) {
		return "", err
	}
	return latestReleaseFile, nil
}
func NewReleaseExists() bool {
	_, err := os.Stat(latestReleaseFile)
	return !os.IsNotExist(err)
}

func GetFnsZipPath(releaseID string) (string, error) {
	fnsFilePath := filepath.Join("releases", releaseID, "fns.zip")
	if _, err := os.Stat(fnsFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("fns.zip not found for release ID %s", releaseID)
	}
	return fnsFilePath, nil
}
