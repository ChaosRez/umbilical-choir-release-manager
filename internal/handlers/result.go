package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/release_manager"
)

func ResultHandler(rm *release_manager.ReleaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resultReq models.ResultRequest

		log.Debugf("Incoming result: %v", r.Body)

		if err := json.NewDecoder(r.Body).Decode(&resultReq); err != nil {
			log.Errorf("Error decoding result: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Infof("Received result (%s) from ChildID: %s (release id: %v, stage: %s)",
			resultReq.StageSummaries[0].Status, resultReq.ChildID, resultReq.ReleaseID, resultReq.StageSummaries[0].StageName)

		if resultReq.StageSummaries == nil || len(resultReq.StageSummaries) == 0 {
			log.Errorf("StageSummaries is missing or empty in the request")
			http.Error(w, "StageSummaries is required", http.StatusBadRequest)
			return
		}

		if err := rm.StagesTracker.StoreResult(resultReq); err != nil {
			log.Errorf("Error storing result: %v", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		// read your write. TODO: for now, only one (first) stage is expected to be sent
		_, exists := rm.StagesTracker.GetResult(resultReq.ReleaseID, resultReq.StageSummaries[0].StageName, resultReq.ChildID)
		if !exists {
			log.Errorf("Error getting result from db (read-your-write) for the pair: %s : %s", resultReq.ChildID, resultReq.ReleaseID)
			http.Error(w, "Error saving the result", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
