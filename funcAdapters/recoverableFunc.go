package funcAdapters

import (
	"github.com/breathbath/go_async_executor/dto"
	"github.com/breathbath/go_async_executor/output"
)

type NamedFunc struct {
	name string
}

func (nf NamedFunc) GetFuncName() string {
	return nf.name
}

/*
RecoverableNonReturningFunc Recoverable means it will try to repeat this function if it fails
see ExecutionSettings.FailedMessagesRepeatAttemptsCount to define how many times
for amqp adapter see ExecutionBuildSettings.FailedMessagesRepeatAttemptsCount for repeat times count and
ExecutionBuildSettings.FailedMessagesRepeatDelay for defining delay time between the attempts
If the function fails more than FailedMessagesRepeatAttemptsCount times
it's payload will be passed to the ExecutionFacade.errorWriter
in case with amqp this is an error queue, containing all failed messages

NonReturning means function has no output for further processing
 */
type RecoverableNonReturningFunc struct {
	NamedFunc
	callback func(input string) error
}

func NewRecoverableNonReturningFunc(name string, callback func(input string) error) *RecoverableNonReturningFunc {
	return &RecoverableNonReturningFunc{NamedFunc: NamedFunc{name: name}, callback: callback}
}

func (rf *RecoverableNonReturningFunc) Process(input string) (dto.OutputMessage, output.ExecutionOutput) {
	err := rf.callback(input)
	if err != nil {
		return dto.StringOutputMessage(""), output.RecoverableError{Err: err}
	}

	return dto.StringOutputMessage(""), output.NonOutputtableSuccess{}
}

/*
RecoverableReturningFunc Recoverable meaning see RecoverableNonReturningFunc
Returning means function's output will be passed further to the output writer of the ExecutionFacade.outputWriter
In case with amqp this is an existing exchange id
which is provided in ExecutionBuildSettings.OutputResultToExistingExchange
 */
type RecoverableReturningFunc struct {
	NamedFunc
	callback func(input string) (string, error)
}

func NewRecoverableReturningFunc(name string, callback func(input string) (string, error)) *RecoverableReturningFunc {
	return &RecoverableReturningFunc{NamedFunc: NamedFunc{name: name}, callback: callback}
}

func (rf *RecoverableReturningFunc) Process(input string) (dto.OutputMessage, output.ExecutionOutput) {
	strOut, err := rf.callback(input)
	if err != nil {
		return dto.StringOutputMessage(strOut), output.RecoverableError{Err: err}
	}

	return dto.StringOutputMessage(strOut), output.OutputtableSuccess{}
}