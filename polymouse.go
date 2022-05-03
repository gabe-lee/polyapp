package polyapp

type MouseInterface interface {
	GetMouseButtonState(button MouseButton) InputState
	GetMousePosition() Vec2
	SetCallbackOnMouseWheelScroll(op func(offset Vec2))
	SetCallbackOnMouseMove(op func(pos Vec2))
	SetCallbackOnMouseButton(op func(button MouseButton, state InputState))
}

var _ MouseInterface = (*MouseProvider)(nil)

type MouseProvider struct {
	MouseInterface
}

type MouseButton uint8

const (
	Mouse1 MouseButton = iota
	Mouse2
	Mouse3
	Mouse4
	Mouse5
	Mouse6
	Mouse7
	Mouse8
	Mouse9
	Mouse10
	Mouse11
	Mouse12
	Mouse13
	Mouse14
	Mouse15
	Mouse16
	Mouse17
	Mouse18
	Mouse19
	Mouse20
	MouseWheelUp
	MouseWheelDown
)
