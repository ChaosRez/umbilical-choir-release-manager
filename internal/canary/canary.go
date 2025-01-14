package canary

import (
	log "github.com/sirupsen/logrus"
	"time"
	"umbilical-choir-release-master/internal/models"
	RM "umbilical-choir-release-master/internal/release_manager"
	"umbilical-choir-release-master/internal/storage"
)

func WaitForEnoughChildren(rm *RM.ReleaseManager, minClients int) {
	log.Infof("Waiting for at least %d children to start canary", minClients)
	for {
		time.Sleep(5 * time.Second)
		if len(rm.Children) >= minClients {
			log.Infof("Found %d children, ready to start canary", len(rm.Children))
			return
		}
	}
}

// location_sequential canary: %10 to %100 of a location (from release strategy), then next location (local first),
// The parent choose a child and ask for gradual 1 to 100, when child is done, will tell the parent
func RunLocSeqAllChild(rm *RM.ReleaseManager, canaryRelease storage.Release) {
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
	for _, child := range rm.Children {
		for {
			releaseStatus, exists := rm.Releases.GetChildStatus(canaryRelease.ID, child.ID)
			if exists {
				if releaseStatus == models.Done {
					break
				}
			} else {
				log.Fatalf("Release '%s' status not found for child %s", canaryRelease.ID, child.ID)
			}
		}
	}
	log.Infof("All children have finished the canary release")
	println(rm.VisualizeReleases())
}

// global_incremental canary: %10 of all locations, then %20 of all locations (global first),
// if the requirements met and it is a failure/error will return, but if it is a success, will send a preliminary result and waits for the parent's signal
func RunGlobalIncAllChild(rm *RM.ReleaseManager, canaryRelease storage.Release) {
	prevStageName := ""
	for _, stageName := range canaryRelease.StageNames {
		log.Infof("Telling children to start the stage '%s' one after another", stageName)
		for _, child := range rm.Children { //FIXME it is not safe as new clients can join and also leave, take a instance of children at the beginning?
			if stageName == canaryRelease.StageNames[0] { // at the first iter just register the child
				rm.RegisterChildForRelease(child.ID, &canaryRelease)
			} else {
				// Tell the child to finish the previous stage
				rm.MarkStageAsShouldEnd(canaryRelease.ID, prevStageName, child.ID)
			}

			// wait for the child to finish the current stage
			for {
				stageStatus, exists := rm.StagesTracker.GetStatus(canaryRelease.ID, stageName, child.ID)
				if exists {
					if stageStatus == models.SuccessWaiting {
						log.Debugf("Child %s successfully finished the stage '%s' for release %s", child.ID, stageName, canaryRelease.ID)
						break
					} else if stageStatus == models.Failure || stageStatus == models.Error {
						//TODO: rollback, and tell all children to finish/rollback
						log.Fatalf("Child %s failed (%s) on the stage '%s' of release %s", child.ID, stageStatus.String(), stageName, canaryRelease.ID)
					} else if stageStatus == models.Completed {
						log.Errorf("Unexpected! Child %s already finished the stage '%s' with status '%s' for release %s, While it had to wait for the parent before finishing", child.ID, stageName, stageStatus.String(), canaryRelease.ID)
					}
				} else {
					log.Fatalf("Stage status not found for child %s on stage %s", child.ID, stageName)
				}
				time.Sleep(1 * time.Second)
			}
		}
		log.Infof("All children finished the stage '%s'", stageName)
		prevStageName = stageName
	}
	log.Infof("All canary stages are finished for release %s, waiting for the final results", canaryRelease.ID)
	for _, child := range rm.Children { // mark clients as should end (WaitForSignal)
		currentStatus, exists := rm.StagesTracker.GetStatus(canaryRelease.ID, prevStageName, child.ID)
		if !exists {
			log.Fatalf("Stage status not found for child %s on last stage %s", child.ID, prevStageName)
		}
		if currentStatus == models.SuccessWaiting {
			rm.MarkStageAsShouldEnd(canaryRelease.ID, prevStageName, child.ID)
		} else {
			log.Warnf("Expected to mark child %s as should end, but it is not in 'SuccessWaiting' state, instead '%s'", child.ID, currentStatus)
		}
	}

	// Wait for all clients to finish the last stage
	for _, child := range rm.Children {
		for {
			stageStatus, exists := rm.StagesTracker.GetStatus(canaryRelease.ID, prevStageName, child.ID)
			if exists {
				if stageStatus != models.ShouldEnd { // => stageStatus == models.Completed || stageStatus == models.Failure || stageStatus == models.Error
					log.Infof("Child %s finished the last stage '%s' with status '%s' for release %s", child.ID, prevStageName, stageStatus.String(), canaryRelease.ID)
					break
				}
			} else {
				log.Fatalf("Stage status not found for child %s on last stage %s", child.ID, prevStageName)
			}
			time.Sleep(1 * time.Second)
		}
	}
	log.Infof("All clients have finished the last stage '%s' for release %s", prevStageName, canaryRelease.ID)
	println(rm.VisualizeReleases())
}
