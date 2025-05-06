package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func parseEvents(inputFilename, outputFilename string, cfg *Config) error {
	file, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	_, err = os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}

	competitors := make(map[int]*Competitor)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		event, err := parseEventLine(line)
		if err != nil {
			return err
		}
		err = writeOutputLog(outputFilename, event)
		if err != nil {
			return err
		}
		processEvent(outputFilename, event, competitors, cfg)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	err = writeResultingTable(outputFilename, competitors, cfg)
	if err != nil {
		return err
	}

	return nil
}

func parseEventLine(line string) (Event, error) {
	parts := strings.SplitN(line, "]", 2)
	if len(parts) != 2 {
		return Event{}, errors.New("bad format")
	}

	eventTime := strings.TrimSpace(parts[0][1:])
	if _, err := time.Parse("15:04:05.000", eventTime); err != nil {
		return Event{}, err
	}

	fields := strings.Fields(strings.TrimSpace(parts[1]))
	if len(fields) < 2 {
		return Event{}, errors.New("bad format")
	}

	var eventId, competitorId int
	if _, err := fmt.Sscanf(fields[0], "%d", &eventId); err != nil {
		return Event{}, errors.New("bad format")
	}
	if _, err := fmt.Sscanf(fields[1], "%d", &competitorId); err != nil {
		return Event{}, errors.New("bad format")
	}

	extra := fields[2:]
	extraStr := strings.Join(extra, " ")

	return Event{
		Time:         eventTime,
		Id:           eventId,
		CompetitorId: competitorId,
		Extra:        extraStr,
	}, nil
}

func processEvent(outputFilename string, event Event, competitors map[int]*Competitor, config *Config) {
	c, exists := competitors[event.CompetitorId]
	eventTime, _ := time.Parse("15:04:05.000", event.Time)

	switch event.Id {
	case 1:
		if exists {
			return
		}
		competitors[event.CompetitorId] = &Competitor{
			Id: event.CompetitorId,
		}
	case 2:
		eventTime, _ := time.Parse("15:04:05.000", event.Extra)
		c.ScheduledStart = eventTime
	case 3:
		// ignored
	case 4:
		c.currentStart = c.ScheduledStart
		if eventTime.Before(c.ScheduledStart) ||
			eventTime.After(c.ScheduledStart.Add(config.parsedStartDelta)) {
			c.Status = "disqualified"
			err := writeOutputLog(outputFilename, Event{
				Time:         event.Time,
				Id:           32,
				CompetitorId: event.CompetitorId,
				Extra:        "missed the start time",
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	case 5:
		// ignored
	case 6:
		c.Hits++
	case 7:
		c.PenaltyLaps += 5*(len(c.Laps)+1) - c.Hits
	case 8:
		c.currentPenaltyStart = eventTime
	case 9:
		penaltyTime := eventTime.Sub(c.currentPenaltyStart)
		c.Penalties += penaltyTime
	case 10:
		lapTime := eventTime.Sub(c.currentStart)
		c.Laps = append(c.Laps, lapTime)
		if len(c.Laps) == config.Laps {
			c.Status = "finished"
			c.Finish = eventTime
			err := writeOutputLog(outputFilename, Event{
				Time:         event.Time,
				Id:           33,
				CompetitorId: event.CompetitorId,
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			c.currentStart = eventTime
		}
	case 11:
		c.Extra = event.Extra
		err := writeOutputLog(outputFilename, Event{
			Time:         event.Time,
			Id:           32,
			CompetitorId: event.CompetitorId,
			Extra:        c.Extra,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
