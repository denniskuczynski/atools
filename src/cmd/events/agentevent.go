package main

import (
	"regexp"
	"time"
)

type AgentEventType int

const (
	None = iota
	GoalStateEvent
	DirectorPlanEvent
	DirectorStepEvent
)

type AgentEvent struct {
	Type      AgentEventType
	EventTime time.Time
	Director  string
	Plan      string
	Move      string
	Step      string
}

func convertToEvent(logLine LogLine) (AgentEvent, error) {
	message := logLine.Message

	// check if goal state
	goalStateRE := regexp.MustCompile("^All (\\d+) Mongo processes are in goal state")

	goalStateMatches := goalStateRE.FindStringSubmatch(message)
	if goalStateMatches != nil {
		return AgentEvent{GoalStateEvent, logLine.LogTime, "", "", "", ""}, nil
	}

	// check if director plan event
	directorPlanRE := regexp.MustCompile("^... process has a plan : (\\S+)")

	directorPlanMatches := directorPlanRE.FindStringSubmatch(message)
	if directorPlanMatches != nil {
		plan := directorPlanMatches[1]
		return AgentEvent{DirectorPlanEvent, logLine.LogTime, logLine.Director, plan, "", ""}, nil
	}

	// check if director step event
	directorStepRE := regexp.MustCompile("^Running step '(\\S+)' as part of move '(\\S+)'")

	directorStepMatches := directorStepRE.FindAllStringSubmatch(message, -1)
	if directorStepMatches != nil {
		step := directorStepMatches[0][1]
		move := directorStepMatches[0][2]
		return AgentEvent{DirectorStepEvent, logLine.LogTime, logLine.Director, "", move, step}, nil
	}

	return AgentEvent{None, time.Time{}, "", "", "", ""}, nil
}
