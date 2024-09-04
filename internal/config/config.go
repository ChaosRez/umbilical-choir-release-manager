package config

import (
	"encoding/json"
	"errors"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	IPAddress          string `yaml:"ip_address"`
	Port               string `yaml:"port"`
	Loglevel           string `yaml:"log_level"`
	ServiceArea        string `yaml:"service_area"`
	ServiceAreaPolygon orb.Polygon
}

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateIPAddress(config.IPAddress); err != nil {
		return nil, err
	}
	if err := validatePort(config.Port); err != nil {
		return nil, err
	}
	config.ServiceAreaPolygon, err = parseServiceAreaPolygon(config.ServiceArea)
	if err != nil {
		return nil, err
	}

	log.Info("Successfully read the config file: ", filename)
	return &config, nil
}

// private
func validateIPAddress(ip string) error {
	if ip == "" {
		return errors.New("IP address is missing or empty")
	}
	return nil
}

func validatePort(port string) error {
	if port == "" {
		return errors.New("Port is missing or empty")
	}
	return nil
}
func parseServiceAreaPolygon(serviceArea string) (orb.Polygon, error) {
	log.Debug("Parsing service area")
	var fc geojson.FeatureCollection
	err := json.Unmarshal([]byte(serviceArea), &fc)
	if err != nil {
		return nil, err
	}

	if len(fc.Features) == 0 {
		return nil, errors.New("Service area is empty or invalid")
	}

	polygon, ok := fc.Features[0].Geometry.(orb.Polygon)
	if !ok {
		return nil, errors.New("Service area is not a valid polygon")
	}

	return polygon, nil
}
