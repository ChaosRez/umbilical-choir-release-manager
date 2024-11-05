package models

import (
	"github.com/paulmach/orb"
	"time"
)

type Parent struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Child struct {
	ID             string      `json:"id"`
	GeographicArea orb.Polygon `json:"geographic_area"`
	LastPoll       time.Time   `json:"last_poll"`
	// TODO: add number of children
}
