package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/models"
)

func ResultHandler(rm *models.ReleaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resultReq models.ResultRequest

		log.Debugf("Incoming result: %v", r.Body)

		if err := json.NewDecoder(r.Body).Decode(&resultReq); err != nil {
			log.Errorf("Error decoding result: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Infof("Received result from child ID: %s", resultReq.ID)
		log.Debugf("Result summary: %+v", resultReq.ReleaseSummary)

		// TODO process the result summary, e.g., store it

		w.WriteHeader(http.StatusOK)
	}
}
