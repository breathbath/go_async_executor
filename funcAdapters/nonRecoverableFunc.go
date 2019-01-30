package funcAdapters

import (
	"async_executor/dto"
	"async_executor/output"
)


/*
NonRecoverableNonReturningFunc NonRecoverable means if function will fail it will not be repeated and the
payload will be passed to the ExecutionFacade.errorWriter
in case with amqp this is an error queue, containing all failed messages

NonReturning see RecoverableNonReturningFunc for explanation
 */
type NonRecoverableNonReturningFunc struct {
	RecoverableNonReturningFunc
}

func NewNonRecoverableNonReturningFunc(name string, callback func(input string) error) *NonRecoverableNonReturningFunc {
	baseFunc := NewRecoverableNonReturningFunc(name, callback)
	return &NonRecoverableNonReturningFunc{RecoverableNonReturningFunc: *baseFunc}
}

func (rf *NonRecoverableNonReturningFunc) Process(input string) (dto.OutputMessage, output.ExecutionOutput) {
	err := rf.callback(input)
	if err != nil {
		return dto.StringOutputMessage(""), output.NonRecoverableError{Err: err}
	}

	return dto.StringOutputMessage(""), output.NonOutputtableSuccess{}
}

/*
NonRecoverableNonReturningFunc for NonRecoverable see NonRecoverableNonReturningFunc for explanation
Returning see RecoverableReturningFunc for explanation
 */
type NonRecoverableReturningFunc struct {
	RecoverableReturningFunc
}

func NewNonRecoverableReturningFunc(name string, callback func(input string) (string, error)) *NonRecoverableReturningFunc {
	baseFunc := NewRecoverableReturningFunc(name, callback)
	return &NonRecoverableReturningFunc{RecoverableReturningFunc: *baseFunc}
}

func (rf *NonRecoverableReturningFunc) Process(input string) (dto.OutputMessage, output.ExecutionOutput) {
	strOutput, err := rf.callback(input)
	if err != nil {
		return dto.StringOutputMessage(strOutput), output.NonRecoverableError{Err: err}
	}

	return dto.StringOutputMessage(strOutput), output.OutputtableSuccess{}
}
