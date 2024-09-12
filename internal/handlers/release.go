package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strings"
	"sync/atomic"

	"umbilical-choir-release-master/internal/repository"
)

var releaseHandlerCounter uint64
var functionsHandlerCounter uint64

// ReleaseHandler serves the latest release.yml file
func ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&releaseHandlerCounter, 1)
	log.Debugf("ReleaseHandler called %d times", atomic.LoadUint64(&releaseHandlerCounter))

	releaseFile, err := repository.GetLatestRelease()
	if err != nil {
		log.Errorf("Error getting latest release: %v", err)
		http.Error(w, "No release found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath.Join(".", releaseFile))
}

// FunctionsHandler serves the fns.zip file for a given release ID
func FunctionsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&functionsHandlerCounter, 1)
	log.Debugf("FunctionsHandler called %d times", atomic.LoadUint64(&functionsHandlerCounter))

	// Extract release_id from URL
	releaseID := strings.TrimPrefix(r.URL.Path, "/release/functions/")
	if releaseID == "" {
		http.Error(w, "Release ID not specified", http.StatusBadRequest)
		return
	}

	// Get the path to the fns.zip file using the repository package
	fnsFilePath, err := repository.GetFnsZipPath(releaseID)
	if err != nil {
		log.Errorf("Error getting fns.zip path for release ID %s: %v", releaseID, err)
		http.Error(w, "related functions not found", http.StatusNotFound)
		return
	}

	// Serve the fns.zip file
	http.ServeFile(w, r, fnsFilePath)
}
