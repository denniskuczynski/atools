package main

import (
	"testing"

	"fmt"
	"time"
)

func TestExtractLogLines(t *testing.T) {
	var tests = []struct {
		logFilePath   string
		expectedCount int
	}{
		{"emptyConfig.log", 3},
		{"addStandalone.log", 15},
		{"addReplicaSet.log", 31},
		{"addCluster.log", 133},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("example%d", i), func(t *testing.T) {
			logLines, err := extractLogLines("../../../exampleLogs/" + test.logFilePath)
			if err != nil {
				t.Errorf("Unable to extract log lines: %v\n", err)
			}

			if len(*logLines) != test.expectedCount {
				t.Errorf("Actual log lines != expected: %v != %v\n", len(*logLines), test.expectedCount)
			}
		})
	}
}

func TestExpandMissingLogTimes(t *testing.T) {
	logLines := []LogLine{
		LogLine{
			time.Date(0, 1, 1, 15, 50, 15, 812000000, time.UTC),
			"",
			"",
			"",
			0,
			"",
			"",
		},
		LogLine{
			time.Date(2016, 10, 15, 15, 51, 15, 812000000, time.UTC),
			"",
			"",
			"",
			0,
			"",
			"",
		},
		LogLine{
			time.Date(0, 1, 1, 15, 52, 15, 812000000, time.UTC),
			"",
			"",
			"",
			0,
			"",
			"",
		},
		LogLine{
			time.Date(0, 1, 1, 15, 53, 15, 812000000, time.UTC),
			"",
			"",
			"",
			0,
			"",
			"",
		},
	}

	expandMissingLogTimes(&logLines)

	for _, line := range logLines {
		if line.LogTime.Year() == 0 {
			t.Errorf("Expected all log lines to have year expanded: %v", line)
		}
	}
}

func TestExtractEvents(t *testing.T) {
	var tests = []struct {
		logFilePath   string
		expectedCount int
	}{
		{"emptyConfig.log", 1},
		{"addStandalone.log", 6},
		{"addReplicaSet.log", 18},
		{"addCluster.log", 100},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("example%d", i), func(t *testing.T) {
			logLines, err := extractLogLines("../../../exampleLogs/" + test.logFilePath)
			if err != nil {
				t.Errorf("Unable to extract log lines: %v\n", err)
			}

			expandMissingLogTimes(logLines)

			events, err := extractEvents(logLines)
			if err != nil {
				t.Errorf("Unable to extract events: %v\n", err)
			}

			if len(*events) != test.expectedCount {
				t.Errorf("Actual events != expected: %v != %v\n", len(*events), test.expectedCount)
			}
		})
	}
}
