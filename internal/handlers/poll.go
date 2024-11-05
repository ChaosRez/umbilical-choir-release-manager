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
	"umbilical-choir-release-master/internal/release_manager"
)

type PollRequest struct {
	ID               string           `json:"id"`
	GeographicArea   geojson.Geometry `json:"geographic_area"`
	NumberOfChildren int              `json:"number_of_children"`
}

type PollResponse struct {
	ID         string `json:"id"`
	NewRelease string `json:"new_release"`
}

func PollHandler(rm *release_manager.ReleaseManager) http.HandlerFunc {
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

		// add the child to the release manager, if no id is passed
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

		// Check if release.yml exists
		newRelease := ""
		releaseID, anyPendingReleaseForChild := rm.Releases.GetNextReleaseForChild(pollReq.ID) // FIXME: the child can join after a release and won't be considered for release
		if anyPendingReleaseForChild {
			newRelease = releaseID
		} else {
			log.Debugf("No pending release found for this Child ID %s", pollReq.ID)
		}

		pollResp := PollResponse{
			ID:         pollReq.ID,
			NewRelease: newRelease,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pollResp); err != nil {
			log.Errorf("Error encoding response: %v", err)
		}
	}
}

func parseGeographicArea(geometry geojson.Geometry) (orb.Polygon, error) {
	polygon, ok := geometry.Geometry().(orb.Polygon)
	if !ok {
		return nil, errors.New("geographic area is not a valid polygon")
	}
	return polygon, nil
}
