package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogLine struct {
	LogTime      time.Time
	LogLevel     string
	FileName     string
	FunctionName string
	LineNumber   int
	Director     string
	Message      string
}

func parseLogLine(logLineString string) (LogLine, error) {
	logLine := LogLine{}

	// extract LogTime

	// "[2016/10/19 15:49:44.540]"
	startsWithDateFormatRE := regexp.MustCompile("^\\[(\\d{4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}.\\d{1,3})\\]")

	// "[15:49:44.540]"
	innerDateFormatRE := regexp.MustCompile("\\[(\\d{1,2}:\\d{1,2}:\\d{1,2}.\\d{1,3})\\]")

	matches := startsWithDateFormatRE.FindStringSubmatch(logLineString)
	if matches != nil {
		parsed, err := time.Parse("2006/01/02 15:04:05.000", string(matches[1]))
		if err != nil {
			return logLine, err
		}
		logLine.LogTime = parsed
	} else {
		matches = innerDateFormatRE.FindStringSubmatch(logLineString)
		if matches != nil {
			parsed, err := time.Parse("15:04:05.000", string(matches[1]))
			if err != nil {
				return logLine, err
			}
			logLine.LogTime = parsed
		}
	}

	// extract LogLevel

	// [header.info]
	logLevelRE := regexp.MustCompile("\\[\\w*\\.(\\w+)\\]")

	matches = logLevelRE.FindStringSubmatch(logLineString)
	if matches != nil {
		logLine.LogLevel = matches[1]
	}

	// extract Director

	// <director>
	directorRE := regexp.MustCompile("<(\\S+)\\>")

	matches = directorRE.FindStringSubmatch(logLineString)
	if matches != nil {
		logLine.Director = matches[1]
	}

	// extract FileName/FunctionName/LineNumber

	// [cm/util/sysdep_unix.go:LockAutomationLockFile:170]
	fileInfoRE := regexp.MustCompile("\\[(\\S+.go:\\w+:\\d+)\\]")

	matches = fileInfoRE.FindStringSubmatch(logLineString)
	if matches != nil {
		splitFileInfo := strings.Split(matches[1], ":")
		logLine.FileName = splitFileInfo[0]
		logLine.FunctionName = splitFileInfo[1]
		lineNumber, err := strconv.Atoi(splitFileInfo[2])
		if err != nil {
			return logLine, fmt.Errorf("Error parsing line number for log line: %v, err: %v", logLineString, err)
		}
		logLine.LineNumber = lineNumber
	}

	// extract Message
	lastClosingSquareBracket := strings.LastIndex(logLineString, "]")
	if lastClosingSquareBracket != -1 {
		logLine.Message = strings.TrimSpace(logLineString[(lastClosingSquareBracket + 1):])
	} else {
		logLine.Message = logLineString
	}

	return logLine, nil
}
