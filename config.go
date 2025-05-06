package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type Config struct {
	Laps             int    `json:"laps"`
	LapLen           int    `json:"lapLen"`
	PenaltyLen       int    `json:"penaltyLen"`
	FiringLines      int    `json:"firingLines"`
	Start            string `json:"start"`
	StartDelta       string `json:"startDelta"`
	parsedStart      time.Time
	parsedStartDelta time.Duration
}

func New(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Laps <= 0 {
		return nil, errors.New("laps must be greater than zero")
	}

	if config.LapLen <= 0 {
		return nil, errors.New("lapLen must be greater than zero")
	}

	if config.FiringLines <= 0 {
		return nil, errors.New("firingLines must be greater than zero")
	}

	if config.PenaltyLen <= 0 {
		return nil, errors.New("penaltyLen must be greater than zero")
	}

	startTime, err := time.Parse("15:04:05", config.Start)
	if err != nil {
		return nil, err
	}

	deltaParts := strings.Split(config.StartDelta, ":")
	if len(deltaParts) != 3 {
		return nil, errors.New("invalid start delta format")
	}

	hours, _ := time.ParseDuration(deltaParts[0] + "h")
	minutes, _ := time.ParseDuration(deltaParts[1] + "m")
	seconds, _ := time.ParseDuration(deltaParts[2] + "s")
	startDelta := hours + minutes + seconds

	config.parsedStart = startTime
	config.parsedStartDelta = startDelta

	return &config, nil
}
