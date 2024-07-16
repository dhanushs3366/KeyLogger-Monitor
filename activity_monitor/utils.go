package activitymonitor

import (
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
		return MSC
	case uint16(BTN_TASK):
		return MSC
	default:
		return MSC
	}

}

func UpdateLogFromEventType(eventType Click, activity *LoggedActivity) {

	switch eventType {
	case KEYBOARD_CLICK:
		activity.Key++

	case LT_MOUSE:
		activity.LeftClicks++

	case RT_MOUSE:
		activity.RightClicks++

	case MID_MOUSE:
		activity.MiddleClicks++

	case EXT_MOUSE_1:
		activity.ExtraClicks++
	case EXT_MOUSE_2:
		activity.ExtraClicks++
	case MSC:
		activity.ExtraClicks++
	default:
		activity.ExtraClicks++
	}

}
