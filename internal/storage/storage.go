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
	if _, exists := rs.storage[key]; exists {
		return fmt.Errorf("result for ChildID %s and ReleaseID %s already exists", result.ChildID, result.ReleaseID)
	}

	rs.storage[key] = result
	log.Info("Stored result for ChildID: ", result.ChildID)
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
