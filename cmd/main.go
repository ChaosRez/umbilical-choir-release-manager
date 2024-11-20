package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"umbilical-choir-release-master/internal/config"
	"umbilical-choir-release-master/internal/handlers"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/release_manager"
	"umbilical-choir-release-master/internal/storage"
)

var conf *config.Config
var rm *release_manager.ReleaseManager
var mainRelease storage.Release
var canaryRelease storage.Release

func main() {
	// >>> simulate a canary <<<<
	// location_sequential: %10 to %100 of a location, then next location (local first),
	// The parent choose a child and ask for gradual 1 to 100, when child is done, will tell the parent

	log.Info("Waiting for at least 2 children to start canary")
	for {
		time.Sleep(5 * time.Second)
		if len(rm.Children) > 1 {
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

	//time.Sleep(15 * time.Second)
	//rm.StagesTracker.UpdateStatus(mainRelease.ID, mainRelease.StageNames[0], rm.Children[0].ID, models.ShouldEnd)
	time.Sleep(2 * time.Second)
}

func init() {
	var err error
	conf, err = config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Log initialization
	config.InitLogger(conf.Loglevel)

	// instantiate release manager
	rm = &release_manager.ReleaseManager{
		Host: conf.Host,
		Port: conf.Port,
		Parent: &models.Parent{
			Host: conf.Parent.Host,
			Port: conf.Parent.Port,
		},
		Children:       []*models.Child{},
		GeographicArea: conf.ServiceAreaPolygon,    // FIXME: if not a leaf, set this to union of children
		StagesTracker:  storage.NewStagesTracker(), // details of a strategy (release)
		Releases:       storage.NewReleases(),
	}

	// Register a release
	mainRelease = storage.Release{
		ID:          "21",
		Name:        "ReleaseSieveFunction",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary test sieve", "A/B Test Sieve"},
	}
	canaryRelease = storage.Release{
		ID:          "22",
		Name:        "Canary10To100_LocationSequential",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary sieve 10", "Canary sieve 90"},
	}
	rm.Releases.AddRelease(mainRelease)
	rm.Releases.AddRelease(canaryRelease)

	// Serve handlers in a separate goroutine
	go func() {
		http.HandleFunc("/poll", handlers.PollHandler(rm))
		http.HandleFunc("/release", handlers.ReleaseHandler(rm))
		http.HandleFunc("/release/functions/", handlers.FunctionsHandler)
		http.HandleFunc("/end_stage", handlers.EndStageHandler(rm))
		http.HandleFunc("/result", handlers.ResultHandler(rm))

		log.Infof("running api on port %s", conf.Port)
		if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}
