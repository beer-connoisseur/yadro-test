package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func writeOutputLog(filename string, event Event) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			log.Println(err)
		}
	}(writer)

	_, _ = writer.WriteString("[" + event.Time + "]" + " ")
	switch event.Id {
	case 1:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) registered\n", event.CompetitorId)
	case 2:
		_, _ = fmt.Fprintf(writer, "The start time for the competitor(%v) was set by a draw to %v\n", event.CompetitorId, event.Extra)
	case 3:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) is on the start line\n", event.CompetitorId)
	case 4:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) has started\n", event.CompetitorId)
	case 5:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) is on the firing range(%v)\n", event.CompetitorId, event.Extra)
	case 6:
		_, _ = fmt.Fprintf(writer, "The target(%v) has been hit by competitor(%v)\n", event.Extra, event.CompetitorId)
	case 7:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) left the firing range\n", event.CompetitorId)
	case 8:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) entered the penalty laps\n", event.CompetitorId)
	case 9:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) left the penalty laps\n", event.CompetitorId)
	case 10:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) ended the main lap\n", event.CompetitorId)
	case 11:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) can`t continue: %v\n", event.CompetitorId, event.Extra)
	case 32:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) is disqualified: %v\n", event.CompetitorId, event.Extra)
	case 33:
		_, _ = fmt.Fprintf(writer, "The competitor(%v) finished the race\n", event.CompetitorId)
	}

	return nil
}

func writeResultingTable(filename string, competitors map[int]*Competitor, cfg *Config) error {
	values := make([]*Competitor, 0, len(competitors))
	for _, v := range competitors {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Finish.Sub(values[i].ScheduledStart) < values[j].Finish.Sub(values[j].ScheduledStart)
	})

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			log.Println(err)
		}
	}(writer)

	_, _ = writer.WriteString("\nResulting Table\n")
	for _, value := range values {
		if value.Status == "disqualified" {
			_, _ = writer.WriteString("[NotStarted] ")
		} else if value.Status == "finished" {
			_, _ = fmt.Fprintf(writer, "%v ", formatDuration(value.Finish.Sub(value.ScheduledStart)))
		} else {
			_, _ = writer.WriteString("[NotFinished] ")
		}

		_, _ = fmt.Fprintf(writer, "%v ", value.Id)

		_, _ = writer.WriteString("[")
		for i := 0; i < cfg.Laps; i++ {
			if i < len(value.Laps) {
				v := value.Laps[i]
				if i != cfg.Laps-1 {
					_, _ = fmt.Fprintf(writer, "{%v, %.3f}, ", formatDuration(v), float64(cfg.LapLen)/v.Seconds())
				} else {
					_, _ = fmt.Fprintf(writer, "{%v, %.3f}] ", formatDuration(v), float64(cfg.LapLen)/v.Seconds())
				}
			} else {
				if i != cfg.Laps-1 {
					_, _ = fmt.Fprintf(writer, "{,}, ")
				} else {
					_, _ = fmt.Fprintf(writer, "{,}] ")
				}
			}
		}
		if value.PenaltyLaps != 0 {
			_, _ = fmt.Fprintf(writer, "{%v, %.3f} ", formatDuration(value.Penalties), float64(value.PenaltyLaps*cfg.PenaltyLen)/value.Penalties.Seconds())
		} else {
			_, _ = fmt.Fprintf(writer, "{,} ")
		}
		_, _ = fmt.Fprintf(writer, "%v/%v\n", value.Hits, len(value.Laps)*5)
	}

	return nil
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
