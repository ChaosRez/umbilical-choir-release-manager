package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"sync/atomic"

	"umbilical-choir-release-master/internal/repository"
)

var releaseHandlerCounter uint64

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
