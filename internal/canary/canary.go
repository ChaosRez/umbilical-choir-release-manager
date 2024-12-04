package canary

import (
	log "github.com/sirupsen/logrus"
	"time"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/release_manager"
	"umbilical-choir-release-master/internal/storage"
)

func WaitForChildrenAndStartCanary(rm *release_manager.ReleaseManager, canaryRelease storage.Release, minClients int) {
	for {
		time.Sleep(5 * time.Second)
		if len(rm.Children) >= minClients {
			log.Infof("Starting canary with %d children", len(rm.Children))
			for _, child := range rm.Children {
				rm.RegisterChildForRelease(child.ID, &canaryRelease)

				// wait for nextChild to finish
				for {
					releaseStatus, exists := rm.Releases.GetChildStatus(canaryRelease.ID, child.ID)
					if exists {
						if releaseStatus == models.Done {
							break
						} else if releaseStatus == models.Failed {
							log.Fatalf("Child %s failed (%s) on the release %s", child.ID, releaseStatus.String(), canaryRelease.ID)
						}
					} else {
						log.Fatalf("Release status not found for child %s", child.ID)
					}
					time.Sleep(1 * time.Second)
				}
			}
			break // canary is done. don't repeat the release
		}
	}
}
