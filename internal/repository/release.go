package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// TODO: check if both fns.zip file and release.yml file exist
func ReadRelease(releaseID string) (string, error) {
	releasePath := filepath.Join("releases", releaseID, "release.yml")
	if _, err := os.Stat(releasePath); os.IsNotExist(err) {
		return "", err
	}
	return releasePath, nil
}

func GetFnsZipPath(releaseID string) (string, error) {
	fnsFilePath := filepath.Join("releases", releaseID, "fns.zip")
	if _, err := os.Stat(fnsFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("fns.zip not found for release ID %s", releaseID)
	}
	return fnsFilePath, nil
}
