package polyapp

type TouchInterface interface {
}

var _ TouchInterface = (*TouchProvider)(nil)

type TouchProvider struct {
	TouchInterface
}
