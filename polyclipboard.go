package polyapp

type ClipboardInterface interface {
	SetClipboardText(text string)
	GetClipboardText() string
}

var _ ClipboardInterface = (*ClipboardProvider)(nil)

type ClipboardProvider struct {
	ClipboardInterface
}
