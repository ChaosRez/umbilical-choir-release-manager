package release_manager

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/storage"
)

type ReleaseManager struct {
	Host           string          `yaml:"host"`
	Port           string          `yaml:"port"`
	GeographicArea orb.Polygon     `yaml:"geographic_area"`
	Parent         *models.Parent  `yaml:"parent"`
	Children       []*models.Child `json:"children"`
	Releases       *storage.Releases
	StagesTracker  *storage.Stages
}

func (rm *ReleaseManager) ChildCount() int {
	return len(rm.Children)
}
func (rm *ReleaseManager) AnyChildren() bool {
	return len(rm.Children) > 0
}
func (rm *ReleaseManager) AddChild(child *models.Child) {
	rm.Children = append(rm.Children, child)
	rm.updateGeographicArea()
}
func (rm *ReleaseManager) RegisterChildForRelease(childID string, release *storage.Release) {
	rm.StagesTracker.InitStagesForChild(release.ID, childID, release.StageNames)
	rm.Releases.MarkChildAsTodo(release.ID, childID)
	log.Infof("Registered child %s for release %s. Now, child should receive it", childID, release.ID)
}
func (rm *ReleaseManager) MarkChildAsNotified(releaseID, childID string) {
	release, exists := (*rm.Releases)[releaseID]
	if !exists {
		log.Errorf("release '%s' not found in Releases to be cleared", releaseID)
		return
	}
	if _, exists := release.ChildStatus[childID]; !exists {
		log.Errorf("child '%s' not found for release '%s' to be marked as Doing", childID, releaseID)
		return
	}
	release.ChildStatus[childID] = models.Doing
	(*rm.Releases)[releaseID] = release
	log.Infof("Updated ReleaseStatus to 'Doing' for child %s in release %s", childID, releaseID)

	// Mark the first stage's status as InProgress
	rm.StagesTracker.UpdateStatus(releaseID, release.StageNames[0], childID, models.InProgress)
}
func (rm *ReleaseManager) MarkChildAsFinished(releaseID, childID string, lastStageStatus models.StageStatus) {
	switch lastStageStatus {
	case models.Completed:
		rm.Releases.SetChildStatus(releaseID, childID, models.Done)
	case models.Failure, models.Error:
		rm.Releases.SetChildStatus(releaseID, childID, models.Failed)
		log.Warnf("Child %s failed (%s) on the %s stage", childID, lastStageStatus.String(), releaseID)
	default:
		log.Errorf("unexpected stage status '%s' for child '%s' in release '%s'", lastStageStatus, childID, releaseID)
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
func (rm *ReleaseManager) VisualizeReleases() (string, error) {
	releasesJSON, err := json.MarshalIndent(rm.Releases, "", "  ")
	if err != nil {
		return "", err
	}
	return string(releasesJSON), nil
}
func (rm *ReleaseManager) VisualizeStagesTracker() (string, error) {
	stages := make(map[string]models.StageSummary)
	for key, summary := range *rm.StagesTracker {
		stages[fmt.Sprintf("%s:%s:%s", key.ReleaseID, key.StageName, key.ChildID)] = summary
	}
	stagesJSON, err := json.MarshalIndent(stages, "", "  ")
	if err != nil {
		return "", err
	}
	return string(stagesJSON), nil
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
	log.Infof("Updated geographic area:\n%v", areaJSON)
}
