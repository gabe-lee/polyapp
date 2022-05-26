package polyapp

type App struct {
	Init       func(options any)
	Teardown   func()
	Window     WindowProvider
	Graphics   GraphicsProvider
	Keyboard   KeyboardProvider
	Mouse      MouseProvider
	Touch      TouchProvider
	Controller ControllerProvider
	File       FileProvider
	Audio      AudioProvider
	Clipboard  ClipboardProvider
}
