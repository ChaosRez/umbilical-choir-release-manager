package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/config"
	"umbilical-choir-release-master/internal/handlers"
	"umbilical-choir-release-master/internal/models"
	"umbilical-choir-release-master/internal/storage"
)

var conf *config.Config
var rm *models.ReleaseManager
var rs *storage.ResultStorage

func main() {
	// Serve handlers in a separate goroutine
	go func() {
		http.HandleFunc("/poll", handlers.PollHandler(rm))
		http.HandleFunc("/release", handlers.ReleaseHandler(rm))
		http.HandleFunc("/release/functions/", handlers.FunctionsHandler)
		http.HandleFunc("/result", handlers.ResultHandler(rs))

		log.Infof("running api on port %s", conf.Port)
		if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// TODO: simulate a canary
	// Check for children every 5 seconds
	time.Sleep(5 * time.Second)
	for {
		if len(rm.Children) > 0 {
			log.Infof("First child ID: %s", rm.Children[0].ID)
			rm.MarkChildForNotification("21", rm.Children[0].ID)
			break
		} else {
			log.Infof("No children found, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
	time.Sleep(1000 * time.Second)
}

func init() {
	var err error
	conf, err = config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Log conf
	ll, err := log.ParseLevel(conf.Loglevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "15:04:05.000", FullTimestamp: true})

	// instantiate release manager
	rm = &models.ReleaseManager{
		Host: conf.Host,
		Port: conf.Port,
		Parent: &models.Parent{
			Host: conf.Parent.Host,
			Port: conf.Parent.Port,
		},
		Children:       []*models.Child{},
		GeographicArea: conf.ServiceAreaPolygon, // FIXME: if not leaf, setting this will be union by other children
	}
	// instantiate result storage
	rs = storage.NewResultStorage()
}
