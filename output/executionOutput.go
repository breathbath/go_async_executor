package output

//ExecutionOutput contains info about execution result and what to do with possible errors
type ExecutionOutput interface {
	GetError() error
}

//OutputtableSuccess if returned - execution was successful and we should output the result message to the outputter
type OutputtableSuccess struct {}

func (seo OutputtableSuccess) GetError() error {
	return nil
}

//NonOutputtableSuccess - execution was successful and we should not do anything with the output
type NonOutputtableSuccess struct {}

func (nos NonOutputtableSuccess) GetError() error {
	return nil
}