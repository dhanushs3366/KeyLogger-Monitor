package activitymonitor

type Click uint

const (
	KEYBOARD_CLICK Click = 0
	RT_MOUSE       Click = 1
	LT_MOUSE       Click = 2
	MID_MOUSE      Click = 3
	EXT_MOUSE_1    Click = 4
	EXT_MOUSE_2    Click = 5
	MSC            Click = 6 //miscalleneous keycodes (joystick etc)
)

const (
	KEYBOARD_MIN = 0x1 //keyboard input key codes are reserved from 0x1 to 0x109
	KEYBOARD_MAX = 0x109
)

type MouseClick uint

const (
	BTN_LEFT    MouseClick = 0x110 //Mouse input key codes are reserved from 0x110 to 0x117
	BTN_RIGHT   MouseClick = 0x111
	BTN_MIDDLE  MouseClick = 0x112
	BTN_SIDE    MouseClick = 0x113
	BTN_EXTRA   MouseClick = 0x114
	BTN_FORWARD MouseClick = 0x115
	BTN_BACK    MouseClick = 0x116
	BTN_TASK    MouseClick = 0x117
)

type devices []string

var AllowedDevices = devices{"Asus Keyboard"}

/*
	@TODO:
	1) Add logging for joysticks so to log gaming activities
*/
