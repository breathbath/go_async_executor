package output

type ExecutionOutput interface {
	GetError() error
}

type OutputtableSuccess struct {}

func (seo OutputtableSuccess) GetError() error {
	return nil
}

type NonOutputtableSuccess struct {}

func (nos NonOutputtableSuccess) GetError() error {
	return nil
}