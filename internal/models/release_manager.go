package models

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"time"
)

type Parent struct {
	IPAddress string `json:"ip_address"`
	Port      string `json:"port"`
}

type Child struct {
	ID             string      `json:"id"`
	GeographicArea orb.Polygon `json:"geographic_area"`
	LastPoll       time.Time   `json:"last_poll"`
}

type ReleaseManager struct {
	Parent         *Parent     `json:"parent"`
	Children       []*Child    `json:"children"`
	GeographicArea orb.Polygon `json:"geographic_area"`
}

func (rm *ReleaseManager) ChildNodeCount() int {
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
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(rm.GeographicArea))

	rawJSON, err := json.Marshal(fc)
	if err != nil {
		return "", err
	}

	return string(rawJSON), nil
}

// Private
func (rm *ReleaseManager) updateGeographicArea() {
	// Create a MultiPolygon to hold all polygons
	multiPolygon := orb.MultiPolygon{}

	for _, child := range rm.Children {
		multiPolygon = append(multiPolygon, child.GeographicArea)
	}
	rm.GeographicArea = multiPolygon.Bound().ToPolygon() // TODO: check if it is same as the union of polygons
}
