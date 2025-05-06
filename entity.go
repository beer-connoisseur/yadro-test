package main

import "time"

type Competitor struct {
	Id                  int
	ScheduledStart      time.Time
	Finish              time.Time
	Laps                []time.Duration
	Penalties           time.Duration
	PenaltyLaps         int
	Hits                int
	Status              string
	Extra               string
	currentStart        time.Time
	currentPenaltyStart time.Time
}

type Event struct {
	Time         string
	Id           int
	CompetitorId int
	Extra        string
}
