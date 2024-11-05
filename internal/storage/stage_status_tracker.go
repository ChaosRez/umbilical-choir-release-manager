package storage

import (
	log "github.com/sirupsen/logrus"
	"umbilical-choir-release-master/internal/models"
)

type StageStatusTracker map[models.StageStatusKey]models.StageStatus

func NewStageStatusTracker() *StageStatusTracker {
	sst := make(StageStatusTracker)
	return &sst
}

func (sst StageStatusTracker) InitStagesForChild(releaseID, childID string, stageNames []string) {
	for _, stageName := range stageNames {
		key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
		sst[key] = models.Pending
	}
}

func (sst StageStatusTracker) AddStatus(releaseID, stageName, childID string, status models.StageStatus) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	sst[key] = status
}

func (sst StageStatusTracker) GetStatus(releaseID, stageName, childID string) (models.StageStatus, bool) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	status, exists := sst[key]
	return status, exists
}

func (sst StageStatusTracker) UpdateStatus(strategyID, stageName, childID string, status models.StageStatus) {
	key := models.StageStatusKey{ReleaseID: strategyID, StageName: stageName, ChildID: childID}
	if _, exists := sst[key]; !exists {
		log.Warnf("Key (%s, %s, %s) does not exist. Cannot update status.", strategyID, stageName, childID)
		return
	}
	sst[key] = status
	log.Infof("Updated '%s' status to %s for (%s) %s", stageName, status, strategyID, childID)
}

func (sst StageStatusTracker) DeleteStatus(releaseID, stageName, childID string) {
	key := models.StageStatusKey{ReleaseID: releaseID, StageName: stageName, ChildID: childID}
	delete(sst, key)
}

//func (sst *StageStatusTracker) GetAllChildStatuses(releaseID, stageName string) map[string]models.StageStatus {
//	result := make(map[string]models.StageStatus)
//	for key, status := range sst.Statuses {
//		if key.ReleaseID == releaseID && key.StageName == stageName {
//			result[key.ChildID] = status
//		}
//	}
//	return result
//}
