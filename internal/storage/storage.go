package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"umbilical-choir-release-master/internal/models"
)

type ResultStorage struct {
	mu      sync.RWMutex
	storage map[string]models.ResultRequest
}

func NewResultStorage() *ResultStorage {
	return &ResultStorage{
		storage: make(map[string]models.ResultRequest),
	}
}

func (rs *ResultStorage) StoreResult(result models.ResultRequest) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	key := rs.generateKey(result.ChildID, result.ReleaseID)
	existingResult, exists := rs.storage[key]
	if !exists {
		rs.storage[key] = result
		log.Infof("Release result for the First stage '%v' (%v) for ChildID: %s has been stored", result.ReleaseSummaries[0].StageName, result.ReleaseID, result.ChildID)
		return nil
	}

	for _, newSummary := range result.ReleaseSummaries {
		if !rs.containsSummary(existingResult.ReleaseSummaries, newSummary) {
			existingResult.ReleaseSummaries = append(existingResult.ReleaseSummaries, newSummary)
		} else {
			log.Errorf("Summary for stage '%v' already exists in the result, client '%v'", newSummary.StageName, result.ChildID)
		}
	}

	rs.storage[key] = existingResult // update the result
	var names []string
	for _, summary := range existingResult.ReleaseSummaries {
		names = append(names, summary.StageName)
	}
	log.Infof("Result updated for ChildID: %v (Release ID: %v, Number of received results: %v, Stages: %v)", result.ChildID, result.ReleaseID, len(existingResult.ReleaseSummaries), names)
	return nil
}

func (rs *ResultStorage) GetResult(childID, releaseID string) (models.ResultRequest, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	key := rs.generateKey(childID, releaseID)
	result, exists := rs.storage[key]
	return result, exists
}

// private
// helper function to generate a key for the storage map
func (rs *ResultStorage) generateKey(childID, releaseID string) string {
	return fmt.Sprintf("%s:%s", childID, releaseID)
}

// helper function to check if a summary already exists in the list
func (rs *ResultStorage) containsSummary(summaries []models.ResultSummary, summary models.ResultSummary) bool {
	for _, s := range summaries {
		if s.StageName == summary.StageName {
			return true
		}
	}
	return false
}
