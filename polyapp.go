package polyapp

type LifecycleInterface interface {
	Init(options ...any)
	Teardown()
}

var _ LifecycleInterface = (*LifecycleProvider)(nil)

type LifecycleProvider struct {
	LifecycleInterface
}

type App struct {
	Lifecycle  LifecycleProvider
	Graphics   GraphicsProvider
	Keyboard   KeyboardProvider
	Mouse      MouseProvider
	Touch      TouchProvider
	Controller ControllerProvider
	File       FileProvider
	Audio      AudioProvider
	Clipboard  ClipboardProvider
}
