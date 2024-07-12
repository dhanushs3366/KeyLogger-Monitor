package activitymonitor

import activitylogger "github.com/dhanushs3366/activity-logger"

func IsKeyInputValid(keyInput activitylogger.InputEvent) bool {
	if keyInput.Type == activitylogger.EV_KEY && keyInput.Value == activitylogger.KEY_PRESSED {
		return true
	}
	return false
}
