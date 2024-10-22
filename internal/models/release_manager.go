package models

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"time"
)

type Parent struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Child struct {
	ID             string      `json:"id"`
	GeographicArea orb.Polygon `json:"geographic_area"`
	LastPoll       time.Time   `json:"last_poll"`
	// TODO: add number of children
}

type ReleaseManager struct {
	Host                 string              `yaml:"host"`
	Port                 string              `yaml:"port"`
	GeographicArea       orb.Polygon         `yaml:"geographic_area"`
	Parent               *Parent             `yaml:"parent"`
	Children             []*Child            `json:"children"`
	ReleaseNotifications map[string][]string // track children to notify for a release
}

func (rm *ReleaseManager) ChildCount() int {
	return len(rm.Children)
}
func (rm *ReleaseManager) HasChildren() bool {
	return len(rm.Children) > 0
}
func (rm *ReleaseManager) AddChild(child *Child) {
	rm.Children = append(rm.Children, child)
	rm.updateGeographicArea()
}
func (rm *ReleaseManager) MarkChildForNotification(releaseNumber, childID string) {
	if rm.ReleaseNotifications == nil {
		rm.ReleaseNotifications = make(map[string][]string)
	}
	rm.ReleaseNotifications[releaseNumber] = append(rm.ReleaseNotifications[releaseNumber], childID)
	log.Infof("Marked child %s for release %s", childID, releaseNumber)
}

// GetReleaseForChild returns the (first) release number for a child, if it exists
func (rm *ReleaseManager) GetReleaseForChild(childID string) (string, bool) {
	for releaseNumber, children := range rm.ReleaseNotifications {
		for _, id := range children {
			if id == childID {
				return releaseNumber, true
			}
		}
	}
	return "", false
}
func (rm *ReleaseManager) ClearNotification(releaseNumber, childID string) {
	children, exists := rm.ReleaseNotifications[releaseNumber]
	if !exists {
		log.Errorf("release '%s' not found in ReleaseNotifications to be cleared", releaseNumber)
		return
	}
	childFound := false
	for i, id := range children {
		if id == childID {
			rm.ReleaseNotifications[releaseNumber] = append(children[:i], children[i+1:]...)
			log.Infof("Cleared notification for child %s in release %s", childID, releaseNumber)
			childFound = true
			break
		}
	}
	if !childFound {
		log.Errorf("child '%s' not found for release '%s' to be cleared", childID, releaseNumber)
	} else {
		if len(rm.ReleaseNotifications[releaseNumber]) == 0 { // if no more children to notify, remove the release
			delete(rm.ReleaseNotifications, releaseNumber)
		}
	}
}
func (rm *ReleaseManager) AreaToJSON() (string, error) {
	gj := geojson.NewGeometry(rm.GeographicArea)

	jsonBlob, err := json.Marshal(gj)
	if err != nil {
		return "", err
	}

	return string(jsonBlob), nil
}

// Private
func (rm *ReleaseManager) updateGeographicArea() {
	// Create a MultiPolygon to hold all polygons
	multiPolygon := orb.MultiPolygon{}

	for _, child := range rm.Children {
		multiPolygon = append(multiPolygon, child.GeographicArea)
	}
	rm.GeographicArea = multiPolygon.Bound().ToPolygon() // TODO: validate it is same as the union of polygons
	areaJSON, _ := rm.AreaToJSON()
	log.Debugf("Updated geographic area:\n%v", areaJSON)
}
