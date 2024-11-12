package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"umbilical-choir-release-master/internal/models"
)

// TODO: add mutex lock to operations?
type Stages map[models.StageStatusKey]models.StageSummary

func NewStagesTracker() *Stages {
	sts := make(Stages)
	return &sts
}

func (sts Stages) InitStagesForChild(releaseID, childID string, stageNames []string) {
	for _, stageName := range stageNames { // mark all stages as pending
		key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
		sts[key] = models.StageSummary{StageName: stageName, Status: models.Pending}
	}
}

func (sts Stages) AddStage(releaseID, stageName, childID string, summary models.StageSummary) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	sts[key] = summary
}
func (sts Stages) DeleteStage(releaseID, stageName, childID string) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	delete(sts, key)
}
func (sts Stages) GetResult(releaseID, stageName, childID string) (models.StageSummary, bool) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	result, exists := sts[key]
	return result, exists
}
func (sts Stages) StoreResult(result models.ResultRequest) error { // Note: for now, only one (first) stage is expected to be sent
	key := models.StageStatusKey{ReleaseID: result.ReleaseID, StageName: result.StageSummaries[0].StageName, ChildID: result.ChildID}
	currentState, exists := sts[key]
	if !exists {
		log.Errorf("Error storing result: stage '%v' does not exist in the result (%s), client '%v'. Expected to be initialized", result.StageSummaries[0].StageName, result.ReleaseID, result.ChildID)
		return fmt.Errorf("stage '%v' does not exist in the result (%s), child '%s'. Expected to be initialized", result.StageSummaries[0].StageName, result.ReleaseID, result.ChildID)
	}
	if currentState.Status > models.Pending {
		sts[key] = result.StageSummaries[0] // TODO: store a batch of stages summaries
		log.Infof("Result updated for the %s stage, ChildID: %v (Release ID: %v, Number of received results: %v)", result.StageSummaries[0].StageName, result.ChildID, result.ReleaseID, len(result.StageSummaries))
		return nil
	} else {
		log.Errorf("stage '%v' already exists in the result (%s), client '%v'", result.StageSummaries[0].StageName, result.ReleaseID, result.ChildID)
		return fmt.Errorf("error storing result: stage '%v' already exists in the result (%s), child '%s'", result.StageSummaries[0].StageName, result.ReleaseID, result.ChildID)
	}
}

// >>> functions to CRUD only 'Status' <<<<

func (sts Stages) GetStatus(releaseID, stageName, childID string) (models.StageStatus, bool) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	summary, exists := sts[key]
	return summary.Status, exists
}

func (sts Stages) UpdateStatus(releaseID, stageName, childID string, status models.StageStatus) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	if _, exists := sts[key]; !exists {
		log.Warnf("Key (%s, %s, %s) does not exist. Cannot update status.", releaseID, stageName, childID)
		return
	}
	summary := sts[key]
	summary.Status = status
	sts[key] = summary
	log.Infof("Updated '%s' status to %s for (%s) %s", stageName, status, releaseID, childID)
}

//func (sts *Stages) GetAllChildStatuses(releaseID, stageName string) map[string]models.StageStatus {
//	result := make(map[string]models.StageStatus)
//	for key, status := range sts.Statuses {
//		if key.ReleaseID == releaseID && key.StageName == stageName {
//			result[key.ChildID] = status
//		}
//	}
//	return result
//}

// private
// helper function to check if a summary already exists in a list
func (sts Stages) summaryAlreadyExists(summaries []models.StageSummary, summary models.StageSummary) bool {
	for _, s := range summaries {
		if s.StageName == summary.StageName {
			return true
		}
	}
	return false
}
