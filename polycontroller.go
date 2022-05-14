package polyapp

type ControllerInterface interface {
}

var _ ControllerInterface = (*ControllerProvider)(nil)

type ControllerProvider struct {
	App *App
	ControllerInterface
}
