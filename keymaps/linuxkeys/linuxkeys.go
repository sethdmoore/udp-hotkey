package linuxkeys

var Keys map[string]uint8
var Modifiers map[string]int

func init() {
	/*
		  Source adapted from https://github.com/micmonay/keybd_event/blob/master/keybd_linux.go

			Permission is hereby granted, free of charge, to any person obtaining a copy
			of this software and associated documentation files (the "Software"), to deal
			in the Software without restriction, including without limitation the rights
			to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
			copies of the Software, and to permit persons to whom the Software is
			furnished to do so, subject to the following conditions:

			The above copyright notice and this permission notice shall be included in
			all copies or substantial portions of the Software.

			THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
			IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
			FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
			AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
			LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
			OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
			THE SOFTWARE.
	*/

	Modifiers = map[string]int{
		"_VK_LEFTCTRL":   29,
		"_VK_RIGHTCTRL":  97,
		"_VK_CTRL":       29,
		"_VK_LEFTSHIFT":  42,
		"_VK_RIGHTSHIFT": 54,
		"_VK_SHIFT":      42,
		"_VK_LEFTALT":    56,
		"_VK_RIGHTALT":   100,
		"_VK_ALT":        56,
	}

	Keys = map[string]uint8{
		"VK_RESERVED":         0,
		"VK_ESC":              1,
		"VK_1":                2,
		"VK_2":                3,
		"VK_3":                4,
		"VK_4":                5,
		"VK_5":                6,
		"VK_6":                7,
		"VK_7":                8,
		"VK_8":                9,
		"VK_9":                10,
		"VK_0":                11,
		"VK_MINUS":            12,
		"VK_EQUAL":            13,
		"VK_BACKSPACE":        14,
		"VK_TAB":              15,
		"VK_Q":                16,
		"VK_W":                17,
		"VK_E":                18,
		"VK_R":                19,
		"VK_T":                20,
		"VK_Y":                21,
		"VK_U":                22,
		"VK_I":                23,
		"VK_O":                24,
		"VK_P":                25,
		"VK_LEFTBRACE":        26,
		"VK_RIGHTBRACE":       27,
		"VK_ENTER":            28,
		"VK_A":                30,
		"VK_S":                31,
		"VK_D":                32,
		"VK_F":                33,
		"VK_G":                34,
		"VK_H":                35,
		"VK_J":                36,
		"VK_K":                37,
		"VK_L":                38,
		"VK_SEMICOLON":        39,
		"VK_APOSTROPHE":       40,
		"VK_GRAVE":            41,
		"VK_BACKSLASH":        43,
		"VK_Z":                44,
		"VK_X":                45,
		"VK_C":                46,
		"VK_V":                47,
		"VK_B":                48,
		"VK_N":                49,
		"VK_M":                50,
		"VK_COMMA":            51,
		"VK_DOT":              52,
		"VK_SLASH":            53,
		"VK_KPASTERISK":       55,
		"VK_SPACE":            57,
		"VK_CAPSLOCK":         58,
		"VK_F1":               59,
		"VK_F2":               60,
		"VK_F3":               61,
		"VK_F4":               62,
		"VK_F5":               63,
		"VK_F6":               64,
		"VK_F7":               65,
		"VK_F8":               66,
		"VK_F9":               67,
		"VK_F10":              68,
		"VK_NUMLOCK":          69,
		"VK_SCROLLLOCK":       70,
		"VK_KP7":              71,
		"VK_KP8":              72,
		"VK_KP9":              73,
		"VK_KPMINUS":          74,
		"VK_KP4":              75,
		"VK_KP5":              76,
		"VK_KP6":              77,
		"VK_KPPLUS":           78,
		"VK_KP1":              79,
		"VK_KP2":              80,
		"VK_KP3":              81,
		"VK_KP0":              82,
		"VK_KPDOT":            83,
		"VK_ZENKAKUHANKAKU":   85,
		"VK_102ND":            86,
		"VK_F11":              87,
		"VK_F12":              88,
		"VK_RO":               89,
		"VK_KATAKANA":         90,
		"VK_HIRAGANA":         91,
		"VK_HENKAN":           92,
		"VK_KATAKANAHIRAGANA": 93,
		"VK_MUHENKAN":         94,
		"VK_KPJPCOMMA":        95,
		"VK_KPENTER":          96,
		"VK_KPSLASH":          98,
		"VK_SYSRQ":            99,
		"VK_LINEFEED":         101,
		"VK_HOME":             102,
		"VK_UP":               103,
		"VK_PAGEUP":           104,
		"VK_LEFT":             105,
		"VK_RIGHT":            106,
		"VK_END":              107,
		"VK_DOWN":             108,
		"VK_PAGEDOWN":         109,
		"VK_INSERT":           110,
		"VK_DELETE":           111,
		"VK_MACRO":            112,
		"VK_MUTE":             113,
		"VK_VOLUMEDOWN":       114,
		"VK_VOLUMEUP":         115,
		"VK_POWER":            116, /* SC System Power Down */
		"VK_KPEQUAL":          117,
		"VK_KPPLUSMINUS":      118,
		"VK_PAUSE":            119,
		"VK_SCALE":            120, /* AL Compiz Scale (Expose) */

		"VK_KPCOMMA":   121,
		"VK_HANGEUL":   122,
		"VK_HANGUEL":   122, // VK_HANGEUL
		"VK_HANJA":     123,
		"VK_YEN":       124,
		"VK_LEFTMETA":  125,
		"VK_RIGHTMETA": 126,
		"VK_COMPOSE":   127,

		"VK_STOP":           128, /* AC Stop */
		"VK_AGAIN":          129,
		"VK_PROPS":          130, /* AC Properties */
		"VK_UNDO":           131, /* AC Undo */
		"VK_FRONT":          132,
		"VK_COPY":           133, /* AC Copy */
		"VK_OPEN":           134, /* AC Open */
		"VK_PASTE":          135, /* AC Paste */
		"VK_FIND":           136, /* AC Search */
		"VK_CUT":            137, /* AC Cut */
		"VK_HELP":           138, /* AL Integrated Help Center */
		"VK_MENU":           139, /* Menu (show menu) */
		"VK_CALC":           140, /* AL Calculator */
		"VK_SETUP":          141,
		"VK_SLEEP":          142, /* SC System Sleep */
		"VK_WAKEUP":         143, /* System Wake Up */
		"VK_FILE":           144, /* AL Local Machine Browser */
		"VK_SENDFILE":       145,
		"VK_DELETEFILE":     146,
		"VK_XFER":           147,
		"VK_PROG1":          148,
		"VK_PROG2":          149,
		"VK_WWW":            150, /* AL Internet Browser */
		"VK_MSDOS":          151,
		"VK_COFFEE":         152, /* AL Terminal Lock/Screensaver */
		"VK_SCREENLOCK":     152, /* VK_COFEE */
		"VK_ROTATE_DISPLAY": 153, /* Display orientation for e.g. tablets */
		"VK_DIRECTION":      153, /* VK_ROTATE_DISPLAY */
		"VK_CYCLEWINDOWS":   154,
		"VK_MAIL":           155,
		"VK_BOOKMARKS":      156, /* AC Bookmarks */
		"VK_COMPUTER":       157,
		"VK_BACK":           158, /* AC Back */
		"VK_FORWARD":        159, /* AC Forward */
		"VK_CLOSECD":        160,
		"VK_EJECTCD":        161,
		"VK_EJECTCLOSECD":   162,
		"VK_NEXTSONG":       163,
		"VK_PLAYPAUSE":      164,
		"VK_PREVIOUSSONG":   165,
		"VK_STOPCD":         166,
		"VK_RECORD":         167,
		"VK_REWIND":         168,
		"VK_PHONE":          169, /* Media Select Telephone */
		"VK_ISO":            170,
		"VK_CONFIG":         171, /* AL Consumer Control Configuration */
		"VK_HOMEPAGE":       172, /* AC Home */
		"VK_REFRESH":        173, /* AC Refresh */
		"VK_EXIT":           174, /* AC Exit */
		"VK_MOVE":           175,
		"VK_EDIT":           176,
		"VK_SCROLLUP":       177,
		"VK_SCROLLDOWN":     178,
		"VK_KPLEFTPAREN":    179,
		"VK_KPRIGHTPAREN":   180,
		"VK_NEW":            181, /* AC New */
		"VK_REDO":           182, /* AC Redo/Repeat */

		"VK_F13": 183,
		"VK_F14": 184,
		"VK_F15": 185,
		"VK_F16": 186,
		"VK_F17": 187,
		"VK_F18": 188,
		"VK_F19": 189,
		"VK_F20": 190,
		"VK_F21": 191,
		"VK_F22": 192,
		"VK_F23": 193,
		"VK_F24": 194,

		"VK_PLAYCD":         200,
		"VK_PAUSECD":        201,
		"VK_PROG3":          202,
		"VK_PROG4":          203,
		"VK_DASHBOARD":      204, /* AL Dashboard */
		"VK_SUSPEND":        205,
		"VK_CLOSE":          206, /* AC Close */
		"VK_PLAY":           207,
		"VK_FASTFORWARD":    208,
		"VK_BASSBOOST":      209,
		"VK_PRINT":          210, /* AC Print */
		"VK_HP":             211,
		"VK_CAMERA":         212,
		"VK_SOUND":          213,
		"VK_QUESTION":       214,
		"VK_EMAIL":          215,
		"VK_CHAT":           216,
		"VK_SEARCH":         217,
		"VK_CONNECT":        218,
		"VK_FINANCE":        219, /* AL Checkbook/Finance */
		"VK_SPORT":          220,
		"VK_SHOP":           221,
		"VK_ALTERASE":       222,
		"VK_CANCEL":         223, /* AC Cancel */
		"VK_BRIGHTNESSDOWN": 224,
		"VK_BRIGHTNESSUP":   225,
		"VK_MEDIA":          226,
		/* Cycle between available video,
		   "outputs" (Monitor/LCD/TV-out/etc) */
		"VK_SWITCHVIDEOMODE": 227,
		"VK_KBDILLUMTOGGLE":  228,
		"VK_KBDILLUMDOWN":    229,
		"VK_KBDILLUMUP":      230,

		"VK_SEND":        231, /* AC Send */
		"VK_REPLY":       232, /* AC Reply */
		"VK_FORWARDMAIL": 233, /* AC Forward Msg */
		"VK_SAVE":        234, /* AC Save */
		"VK_DOCUMENTS":   235,

		"VK_BATTERY": 236,

		"VK_BLUETOOTH": 237,
		"VK_WLAN":      238,
		"VK_UWB":       239,

		"VK_UNKNOWN": 240,

		"VK_VIDEO_NEXT": 241, /* drive next video source */
		"VK_VIDEO_PREV": 242, /* drive previous video source */
		/* brightness up, after max is min
		Set Auto Brightness: manual,
		brightness control is off
		rely on ambient */
		"VK_BRIGHTNESS_CYCLE": 243,
		"VK_BRIGHTNESS_AUTO":  244,
		"VK_BRIGHTNESS_ZERO":  244, /* VK_BRIGHTNESS_AUTO */
		"VK_DISPLAY_OFF":      245, /* display device to off state */

		"VK_WWAN":   246, /* Wireless WAN (LTE, UMTS, GSM, etc.) */
		"VK_WIMAX":  246, /* VK_WWAN */
		"VK_RFKILL": 247, /* Key that controls all radios */

		"VK_MICMUTE": 248, /* Mute / unmute the microphone */
	}

	_ = Keys
	_ = Modifiers
}
