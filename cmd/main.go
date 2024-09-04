package main

import (
	log "github.com/sirupsen/logrus"
	"umbilical-choir-release-master/internal/config"
	"umbilical-choir-release-master/internal/models"
)

var conf *config.Config

func main() {
	rm := &models.ReleaseManager{
		Parent: &models.Parent{
			IPAddress: conf.IPAddress,
			Port:      conf.Port,
		},
		Children:       []*models.Child{},
		GeographicArea: conf.ServiceAreaPolygon,
	}
	log.Infof("%v", rm.Children)

	//http.HandleFunc("/releases", handlers.ReleaseHandler)
	//http.ListenAndServe(":8080", nil)
}

func init() {
	var err error
	conf, err = config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	ll, err := log.ParseLevel(conf.Loglevel)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "15:04:05.000", FullTimestamp: true})
}
