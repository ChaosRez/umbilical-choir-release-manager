package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	Canary "umbilical-choir-release-master/internal/canary"
	"umbilical-choir-release-master/internal/config"
	"umbilical-choir-release-master/internal/handlers"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/release_manager"
	"umbilical-choir-release-master/internal/storage"
)

var conf *config.Config
var rm *release_manager.ReleaseManager
var mainRelease storage.Release
var canaryReleaseLocSeq storage.Release
var canaryReleaseGlobInc storage.Release
var canary5percentRelease storage.Release
var abRelease storage.Release

func main() {

	// 1. do canary on all children
	Canary.WaitForEnoughChildren(rm, 2)
	Canary.RunGlobalIncAllChild(rm, canary5percentRelease) // initial canary test on all children
	for _, child := range rm.Children {                    // check if Release is successful
		childStatus, _ := rm.Releases.GetChildStatus(canary5percentRelease.ID, child.ID)
		if childStatus == models.Done {
			// Ok. all children expected to finish + were successful
		} else {
			log.Fatalf("the child '%s' is not done (%v) on the release %s", child.ID, childStatus, canary5percentRelease.ID)
		}
	}

	// 2. A/B on a subset of nodes
	selectedChildren := rm.Children[:1] // take the first child
	for _, child := range selectedChildren {
		rm.RegisterChildForRelease(child.ID, &abRelease)
	}
	for _, child := range selectedChildren { // check if Release is successful
		// wait for nextChild to finish
		for {
			releaseStatus, exists := rm.Releases.GetChildStatus(abRelease.ID, child.ID)
			if exists {
				if releaseStatus == models.Done {
					break
				} else if releaseStatus == models.Failed {
					log.Fatalf("Child %s failed (%s) on the release %s", child.ID, releaseStatus.String(), abRelease.ID)
				} else {
					time.Sleep(1 * time.Second)
					continue // retry when in progress
				}
			} else {
				log.Fatalf("Release status not found for child %s", child.ID)
			}
		}
	}

	// 3. gradual canary
	Canary.RunGlobalIncAllChild(rm, canaryReleaseGlobInc)

	//Canary.RunLocSeqAllChild(rm, canaryReleaseLocSeq)
	//Canary.RunGlobalIncAllChild(rm, canaryReleaseGlobInc)
	time.Sleep(10 * time.Second)
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
	canaryReleaseLocSeq = storage.Release{
		ID:          "22",
		Name:        "Canary10To100_LocationSequential",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary sieve 10", "Canary sieve 90"},
	}
	canaryReleaseGlobInc = storage.Release{
		ID:          "23",
		Name:        "Canary10To100_GlobalIncremental",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary sieve 10", "Canary sieve 50", "Canary sieve 90"},
	}
	canary5percentRelease = storage.Release{
		ID:          "11",
		Name:        "canary5percent",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary 5 Percent"},
	}
	abRelease = storage.Release{
		ID:          "12",
		Name:        "A/BTestSieveFunction",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"A/B Test Sieve"},
	}

	rm.Releases.AddRelease(mainRelease)
	rm.Releases.AddRelease(canaryReleaseLocSeq)
	rm.Releases.AddRelease(canaryReleaseGlobInc)
	rm.Releases.AddRelease(canary5percentRelease)
	rm.Releases.AddRelease(abRelease)

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
