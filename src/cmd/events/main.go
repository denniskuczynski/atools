package main

import (
	"flag"
	"fmt"
	"os"
)

// used for the 'f' and 'config' flags
var filename string

func main() {
	flag.StringVar(&filename, "filename", "", "file path to automation log file")
	flag.Parse()

	if filename == "" {
		fmt.Fprintln(os.Stderr, "-filename must be set")
		fmt.Fprintln(os.Stderr, UsageMsg)
		os.Exit(1)
	}

	logLines, err := extractLogLines(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting log lines: %v\n", err)
		os.Exit(1)
	}

	expandMissingLogTimes(logLines)

	events, err := extractEvents(logLines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting events: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Printing events")
	for _, event := range *events {
		fmt.Println(event)
	}

	os.Exit(0)
}

var UsageMsg = `Usage :     
REQUIRED : 

        - filename <string>      # file path to automation log file
`
