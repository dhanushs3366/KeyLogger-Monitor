package activitymonitor

import (
	"log"

	activitylogger "github.com/dhanushs3366/activity-logger"
)

func IsKeyInputValid(keyInput activitylogger.InputEvent) bool {
	if keyInput.Type == activitylogger.EV_KEY && keyInput.Value == activitylogger.KEY_PRESSED {
		return true
	}
	return false
}

func CategorizeEvent(inputEvent activitylogger.InputEvent) Click {
	keyCode := inputEvent.Code
	if keyCode >= KEYBOARD_MIN && keyCode <= KEYBOARD_MAX {
		return KEYBOARD_CLICK
	}
	switch keyCode {
	case uint16(BTN_LEFT):
		return LT_MOUSE
	case uint16(BTN_RIGHT):
		return RT_MOUSE
	case uint16(BTN_MIDDLE):
		return MID_MOUSE
	case uint16(BTN_SIDE):
		return EXT_MOUSE_1
	case uint16(BTN_FORWARD):
		return EXT_MOUSE_2
	case uint16(BTN_BACK):
	case uint16(BTN_TASK):
	default:
		return MSC
	}
	return MSC
}

func UpdateLogFromEventType(eventType Click, activity *LoggedActivity) {
	log.SetPrefix("From Updating log")
	switch eventType {
	case KEYBOARD_CLICK:
		activity.Key++
		log.Println(activity)

	case LT_MOUSE:
		activity.LeftClicks++
		log.Println(activity)

	case RT_MOUSE:
		activity.RightClicks++
		log.Println(activity)

	case MID_MOUSE:
		activity.MiddleClicks++
		log.Println(activity)

	case EXT_MOUSE_1:
	case EXT_MOUSE_2:
	case MSC:
	default:
		activity.ExtraClicks++
		log.Println(activity)
	}

}
