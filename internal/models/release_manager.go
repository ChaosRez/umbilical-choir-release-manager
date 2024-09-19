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
	Host           string      `yaml:"host"`
	Port           string      `yaml:"port"`
	GeographicArea orb.Polygon `yaml:"geographic_area"` // FIXME: is optional, therefore can cause problem for its parent. can start polling after registering a child
	Parent         *Parent     `yaml:"parent"`
	Children       []*Child    `json:"children"`
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
