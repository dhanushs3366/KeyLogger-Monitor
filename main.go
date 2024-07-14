package main

import (
	"time"

	activitymonitor "github.com/dhanushs3366/KeyLogger-Monitor/activity_monitor"
)

func main() {
	devPaths := []string{
		// "/dev/input/event6",
		// "/dev/input/event10",
		"/dev/input/event14",
		"/dev/input/event17",
		"/dev/input/event9",
	}
	keyloggers, err := activitymonitor.SetupLoggers(devPaths)
	if err != nil {
		panic(err)
	}
	// keylogger, err := activitylogger.New("/dev/input/event14")
	// if err != nil {
	// 	log.Printf("Error:%v\n", err)
	// }
	// activitymonitor.SendDataFromLogger(*keylogger, 5*time.Second)
	for {
		activitymonitor.SendDataFromLoggers(keyloggers, 5*time.Second)

	}

	// keylogger := keyloggers[len(keyloggers)-1]
	// go activitymonitor.SendDataFromLogger(keylogger, 5*time.Second)
	// Prevent the main function from exiting
	// select {}
}
