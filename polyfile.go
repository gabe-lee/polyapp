package polyapp

type FileInterface interface {
	LoadFileBytes(name string) ([]byte, error)
	SaveFileBytes(name string, data []byte) error
}

var _ FileInterface = (*FileProvider)(nil)

type FileProvider struct {
	App *App
	FileInterface
}

func (f FileProvider) LoadFileString(name string) (string, error) {
	bytes, err := f.LoadFileBytes(name)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (f FileProvider) SaveFileString(name string, content string) error {
	bytes := []byte(content)
	return f.SaveFileBytes(name, bytes)
}
