package main

import (
	"bufio"
	"os"
	"time"
)

func extractLogLines(logFileName string) (*[]LogLine, error) {
	logLines := make([]LogLine, 0)

	log, err := os.Open(logFileName)
	if err != nil {
		return &logLines, err
	}

	index := 0
	lineScanner := bufio.NewScanner(log)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		logLine, err := parseLogLine(line)
		if err != nil {
			return &logLines, err
		}

		// if no time found, assume the logLine belonged to the previous line multi-line message
		if logLine.LogTime.IsZero() && index != 0 {
			logLines[index-1].Message = logLine.Message
		} else {
			logLines = append(logLines, logLine)
			index = index + 1
		}
	}

	if err := lineScanner.Err(); err != nil {
		return &logLines, err
	}

	return &logLines, nil
}

func expandMissingLogTimes(logLines *[]LogLine) {
	var currentYear int
	var currentMonth time.Month
	var currentDay int
	index := -1

	// find first known year/month/day
	for _, logLine := range *logLines {
		index = index + 1
		if logLine.LogTime.Year() != 0 {
			currentYear = logLine.LogTime.Year()
			currentMonth = logLine.LogTime.Month()
			currentDay = logLine.LogTime.Day()
			break
		}
	}

	// start filling in
	for index = 0; index < len(*logLines); index = index + 1 {
		logLine := &((*logLines)[index])
		logTime := logLine.LogTime
		if logTime.Year() == 0 {
			logLine.LogTime = time.Date(currentYear, currentMonth, currentDay, logTime.Hour(), logTime.Minute(), logTime.Second(), logTime.Nanosecond(), logTime.Location())
		}

		currentYear = logLine.LogTime.Year()
		currentMonth = logLine.LogTime.Month()
		currentDay = logLine.LogTime.Day()
	}
}

func extractEvents(logLines *[]LogLine) (*[]AgentEvent, error) {
	events := make([]AgentEvent, 0)

	var lastEvent *AgentEvent
	for _, logLine := range *logLines {
		event, err := convertToEvent(logLine)
		if err != nil {
			return &events, err
		}

		if event.Type != None {
			if (event.Type == GoalStateEvent) &&
				(lastEvent != nil) &&
				(lastEvent.Type == GoalStateEvent) {
				continue
			}

			lastEvent = &event
			events = append(events, event)
		}
	}

	return &events, nil
}
