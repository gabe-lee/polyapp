package polyapp

type KeyboardInterface interface {
	GetKeyboardKeyState(key KeyboardKey) InputState
	SetCallbackOnRuneInput(op func(r rune))
	SetCallbackOnKeyPress(op func(key KeyboardKey, state InputAction, mods KeyboardMod))
}

var _ KeyboardInterface = (*KeyboardProvider)(nil)

type KeyboardProvider struct {
	KeyboardInterface
}

type KeyboardKey uint8

const (
	KeyUnknown KeyboardKey = iota
	KeySpace
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyKbMenu
	KeyLeftBracket
	KeyBackSlash
	KeyRightBracket
	KeyGrave
	KeyKp0
	KeyKp1
	KeyKp2
	KeyKp3
	KeyKp4
	KeyKp5
	KeyKp6
	KeyKp7
	KeyKp8
	KeyKp9
	KeyKpDecimal
	KeyKpDivide
	KeyKpMultiply
	KeyKpSubtract
	KeyKpAdd
	KeyKpEnter
	KeyKpEqual
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	KeyZero
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
)

type KeyboardMod uint8

const (
	ModShift KeyboardMod = 1 << iota
	ModControl
	ModAlt
	ModSuper
	ModCapsLock
	ModNumLock
)
