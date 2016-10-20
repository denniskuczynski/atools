package main

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertToEvent(t *testing.T) {
	var tests = []struct {
		logLine  LogLine
		expected AgentEvent
	}{
		{
			LogLine{
				time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC),
				"info",
				"cm/util/sysdep_unix.go",
				"LockAutomationLockFile",
				170,
				"",
				"Locking automation lock file at /tmp/mongodb-mms-automation.lock",
			},
			AgentEvent{None, time.Time{}, "", "", "", ""},
		},
		{
			LogLine{
				time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC),
				"",
				"",
				"",
				0,
				"",
				"All 0 Mongo processes are in goal state, No monitoring agent specified, No backup agent specified",
			},
			AgentEvent{GoalStateEvent, time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC), "", "", "", ""},
		},
		{
			LogLine{
				time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC),
				"",
				"",
				"",
				0,
				"myStandalone_0",
				"Running step 'Extract' as part of move 'Download'",
			},
			AgentEvent{DirectorStepEvent, time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC), "myStandalone_0", "", "Download", "Extract"},
		},
		{
			LogLine{
				time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC),
				"",
				"",
				"",
				0,
				"myStandalone_0",
				"... process has a plan : Start",
			},
			AgentEvent{DirectorPlanEvent, time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC), "myStandalone_0", "Start", "", ""},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("example%d", i), func(t *testing.T) {
			event, err := convertToEvent(test.logLine)
			if err != nil {
				t.Errorf("Unable to convert to event: %v\n", err)
			}

			if event != test.expected {
				t.Errorf("Event %v != %v", event, test.expected)
			}
		})
	}
}
