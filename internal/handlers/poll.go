package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"umbilical-choir-release-master/internal/models"
)

type PollRequest struct {
	ID               string           `json:"id"`
	GeographicArea   geojson.Geometry `json:"geographic_area"`
	NumberOfChildren int              `json:"number_of_children"`
}

type PollResponse struct {
	ID string `json:"id"`
}

func parseGeographicArea(geometry geojson.Geometry) (orb.Polygon, error) {
	polygon, ok := geometry.Geometry().(orb.Polygon)
	if !ok {
		return nil, errors.New("Geographic area is not a valid polygon")
	}
	return polygon, nil
}

func PollHandler(rm *models.ReleaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pollReq PollRequest

		//log.Debugf("Incoming request: %v", r.Body)

		if err := json.NewDecoder(r.Body).Decode(&pollReq); err != nil {
			log.Errorf("Error decoding request: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Debugf("Decoded request: %+v", pollReq)

		geoArea, err := parseGeographicArea(pollReq.GeographicArea)
		if err != nil {
			log.Errorf("Error parsing geographic area: %v", err)
			http.Error(w, "Invalid geographic area", http.StatusBadRequest)
			return
		}

		if pollReq.ID == "" {
			pollReq.ID = uuid.New().String()
			newChild := &models.Child{
				ID:             pollReq.ID,
				GeographicArea: geoArea,
				LastPoll:       time.Now(),
			}
			rm.AddChild(newChild)
			log.Info("updated child count: ", rm.ChildCount())
		}

		pollResp := PollResponse{ID: pollReq.ID}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pollResp); err != nil {
			log.Errorf("Error encoding response: %v", err)
		}
	}
}
