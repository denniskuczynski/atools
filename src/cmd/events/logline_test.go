package main

import (
	"fmt"
	"testing"
	"time"
)

func TestParseLogLine(t *testing.T) {
	var tests = []struct {
		logLineString string
		expected      LogLine
	}{
		{
			"[2016/10/19 15:49:44.540] [.info] [cm/util/sysdep_unix.go:LockAutomationLockFile:170] [15:49:44.540] Locking automation lock file at /tmp/mongodb-mms-automation.lock",
			LogLine{
				time.Date(2016, time.October, 19, 15, 49, 44, 540000000, time.UTC),
				"info",
				"cm/util/sysdep_unix.go",
				"LockAutomationLockFile",
				170,
				"",
				"Locking automation lock file at /tmp/mongodb-mms-automation.lock",
			},
		},
		{
			"[13:34:14.786] All 0 Mongo processes are in goal state, No monitoring agent specified, No backup agent specified",
			LogLine{
				time.Date(0, 1, 1, 13, 34, 14, 786000000, time.UTC),
				"",
				"",
				"",
				0,
				"",
				"All 0 Mongo processes are in goal state, No monitoring agent specified, No backup agent specified",
			},
		},
		{
			"<myStandalone_0> [15:50:15.812] ... process has a plan : Download,Start",
			LogLine{
				time.Date(0, 1, 1, 15, 50, 15, 812000000, time.UTC),
				"",
				"",
				"",
				0,
				"myStandalone_0",
				"... process has a plan : Download,Start",
			},
		},
		{
			"/desired/      Name = s+x_rs1_1",
			LogLine{
				time.Time{},
				"",
				"",
				"",
				0,
				"",
				"/desired/      Name = s+x_rs1_1",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("example%d", i), func(t *testing.T) {
			logLine, err := parseLogLine(test.logLineString)
			if err != nil {
				t.Errorf("Unable to parse log line: %v\n", err)
			}

			if logLine.LogTime != test.expected.LogTime {
				t.Errorf("LogTime %v != %v", logLine.LogTime, test.expected.LogTime)
			}

			if logLine.LogLevel != test.expected.LogLevel {
				t.Errorf("LogLevel '%v' != '%v'", logLine.LogLevel, test.expected.LogLevel)
			}

			if logLine.Director != test.expected.Director {
				t.Errorf("Director '%v' != '%v'", logLine.Director, test.expected.Director)
			}

			if logLine.FileName != test.expected.FileName {
				t.Errorf("FileName '%v' != '%v'", logLine.FileName, test.expected.FileName)
			}

			if logLine.FunctionName != test.expected.FunctionName {
				t.Errorf("FunctionName '%v' != '%v'", logLine.FunctionName, test.expected.FunctionName)
			}

			if logLine.LineNumber != test.expected.LineNumber {
				t.Errorf("LineNumber '%v' != '%v'", logLine.LineNumber, test.expected.LineNumber)
			}

			if logLine.FileName != test.expected.FileName {
				t.Errorf("FileName '%v' != '%v'", logLine.FileName, test.expected.FileName)
			}

			if logLine.Message != test.expected.Message {
				t.Errorf("Message '%v' != '%v'", logLine.Message, test.expected.Message)
			}
		})
	}
}
