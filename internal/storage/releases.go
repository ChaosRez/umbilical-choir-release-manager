package storage

import (
	log "github.com/sirupsen/logrus"
	"umbilical-choir-release-master/internal/models"
)

type Releases map[string]Release

type Release struct {
	ID          string
	Name        string
	Type        string
	Functions   []string
	StageNames  []string
	ChildStatus map[string]models.ReleaseStatus
}

func NewReleases() *Releases {
	releases := make(Releases)
	return &releases
}

func (r Releases) AddRelease(release Release) {
	if _, exists := r[release.ID]; exists {
		log.Errorf("release '%s' already exists in Releases", release.ID)
		return
	}
	r[release.ID] = release
	log.Infof("Added release '%s' to Releases", release.ID)
}

func (r Releases) MarkChildAsTodo(releaseID, childID string) {
	release, exists := r[releaseID]
	if !exists {
		log.Errorf("release '%s' not found in Releases to register the child for it", releaseID)
		return
	}

	release.ChildStatus[childID] = models.Todo // <<<
	r[releaseID] = release
	log.Debugf("Marked child %s for release %s", childID, releaseID)
}

func (r Releases) GetNextReleaseForChild(childID string) (string, bool) {
	for releaseID, release := range r {
		if status, exists := release.ChildStatus[childID]; exists {
			if status == models.Todo { // pending release for the child
				return releaseID, true
			} else {
				log.Debugf("(ReleaseStatus: %v) release %s is not supposed to be run by the child %s", status, releaseID, childID)
			}
		}
	}
	return "", false
}
