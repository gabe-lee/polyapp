package polyapp

type AudioInterface interface {
}

var _ AudioInterface = (*AudioProvider)(nil)

type AudioProvider struct {
	AudioInterface
}
