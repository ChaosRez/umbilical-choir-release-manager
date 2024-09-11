package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"umbilical-choir-release-master/internal/config"
	"umbilical-choir-release-master/internal/handlers"
	"umbilical-choir-release-master/internal/models"
)

var conf *config.Config
var rm *models.ReleaseManager

func main() {
	http.HandleFunc("/poll", handlers.PollHandler(rm))
	http.HandleFunc("/release", handlers.ReleaseHandler)
	http.HandleFunc("/result", handlers.ResultHandler(rm))
	log.Infof("running api on port %s", conf.Port)
	http.ListenAndServe(":"+conf.Port, nil)
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
		GeographicArea: conf.ServiceAreaPolygon,
	}
}
