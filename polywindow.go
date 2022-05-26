package polyapp

import "image"

type WindowInterface interface {
	CreateWindow() (windowID uint8, err error)
	DestroyWindow(windowID uint8) (err error)
	RequestClose(windowID uint8) (err error)
	GetSize(windowID uint8) (size IVec2, err error)
	SetSize(windowID uint8, size IVec2) error
	GetPos(windowID uint8, size IVec2) error
	SetPos(windowID uint8, size IVec2) error
	SetOpacity(windowID uint8, opacity float32) error
	SetTitle(windowID uint8, title string) error
	SetIcon(windowID uint8, icon image.RGBA) error
	SetFocusCallback(windowID uint8, op func(focused bool)) error
	SetCloseCallback(windowID uint8, op func()) error
	SetMinimizeCallback(windowID uint8, op func(minimized bool)) error
	SetMaximizeCallback(windowID uint8, op func(maximized bool)) error
	SetPosCallback(windowID uint8, op func(pos IVec2)) error
	SetSizeCallback(windowID uint8, op func(size IVec2)) error
}

var _ WindowInterface = (*WindowProvider)(nil)

type WindowProvider struct {
	WindowInterface
}
