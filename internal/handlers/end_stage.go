package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/release_manager"
)

type EndStageRequest struct {
	ID         string `json:"id"`
	StrategyID string `json:"strategy_id"`
	StageName  string `json:"stage_name"`
}

type EndStageResponse struct {
	EndStage bool `json:"end_stage"`
}

func EndStageHandler(rm *release_manager.ReleaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EndStageRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("Error decoding request: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check if the child exists
		result, exists := rm.StageStatusTracker.GetStatus(req.StrategyID, req.StageName, req.ID)
		if !exists {
			log.Errorf("No stage status found for (strategy: %s, stage: %s, child: %s)", req.StrategyID, req.StageName, req.ID)
			http.Error(w, "No stage status found for mentioned combination", http.StatusNotFound)
			return
		}

		// Prepare the response
		response := EndStageResponse{
			EndStage: result >= models.ShouldEnd, // if the stage is ShouldEnd or higher, end the stage
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Errorf("Error encoding response: %v", err)
		}
	}
}
