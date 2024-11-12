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

func main() {
	mainRelease := storage.Release{
		ID:          "21",
		Name:        "ReleaseSieveFunction",
		Type:        "major",
		Functions:   []string{"sieve"},
		ChildStatus: map[string]models.ReleaseStatus{},
		StageNames:  []string{"Canary test sieve", "A/B Test Sieve"},
	}
	rm.Releases.AddRelease(mainRelease)
	// TODO: simulate a canary
	// Check for children every 5 seconds
	time.Sleep(5 * time.Second)
	for {
		if len(rm.Children) > 0 {
			log.Infof("First child ID: %s", rm.Children[0].ID)
			//rm.Releases.MarkChildAsTodo("21", rm.Children[0].ID)
			rm.RegisterChildForRelease(rm.Children[0].ID, &mainRelease)
			break
		} else {
			log.Infof("No children found, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
	time.Sleep(15 * time.Second)
	rm.StagesTracker.UpdateStatus(mainRelease.ID, "Canary test sieve", rm.Children[0].ID, models.ShouldEnd)
	time.Sleep(1000 * time.Second)
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
