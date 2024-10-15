package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/storage"
)

func ResultHandler(rs *storage.ResultStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resultReq models.ResultRequest

		log.Debugf("Incoming result: %v", r.Body)

		if err := json.NewDecoder(r.Body).Decode(&resultReq); err != nil {
			log.Errorf("Error decoding result: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Infof("Received result (%s) from ChildID: %s (release id: %v, stage: %s)", resultReq.ReleaseSummaries[0].Status, resultReq.ChildID, resultReq.ReleaseID, resultReq.ReleaseSummaries[0].StageName)

		if resultReq.ReleaseSummaries == nil || len(resultReq.ReleaseSummaries) == 0 {
			log.Errorf("ReleaseSummaries is missing or empty in the request")
			http.Error(w, "ReleaseSummaries is required", http.StatusBadRequest)
			return
		}

		if err := rs.StoreResult(resultReq); err != nil {
			log.Errorf("Error storing result: %v", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		// read your write
		_, exists := rs.GetResult(resultReq.ChildID, resultReq.ReleaseID)
		if !exists {
			log.Errorf("Error getting result from db (read-your-write) for the pair: %s : %s", resultReq.ChildID, resultReq.ReleaseID)
			http.Error(w, "Error saving the result", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
