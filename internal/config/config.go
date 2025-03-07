package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Loglevel           string `yaml:"log_level"`
	Parent             Parent `yaml:"parent"`
	ServiceArea        string `yaml:"service_area"`
	ServiceAreaPolygon orb.Polygon
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
}
type Parent struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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

	if err := validateDefined("Host", config.Host); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}
	if err := validateDefined("Port", config.Port); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}
	if err := validateDefined("Parent Host", config.Parent.Host); err != nil {
		log.Warn("Parent host is missing, assuming this is the mother node")
	}
	if config.ServiceArea != "" {
		config.ServiceAreaPolygon, err = parseServiceAreaPolygon(config.ServiceArea)
		if err != nil {
			return nil, fmt.Errorf("config validation error: %w", err)
		}
	} else {
		config.ServiceAreaPolygon = orb.Polygon{}
		log.Warn("Service area is not defined, using an empty polygon")
	}

	log.Info("Successfully read the config file: ", filename)
	return &config, nil
}

// Set logger's level (from config) and format
func InitLogger(logLevel string) {
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "15:04:05.000",
		FullTimestamp:   false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, file := filepath.Split(f.File)
			return "", fmt.Sprintf(" %s:%d", file, f.Line)
		},
	})
}

// private
func validateDefined(field, str string) error {
	if str == "" {
		return errors.New(field + " is missing or empty")
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
