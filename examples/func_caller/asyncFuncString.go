package func_caller

import (
	"go_async_executor/dto"
	"go_async_executor/output"
)

type StringArgFunc func(input string) error

type StringFunctionExecutor struct {
	name string
	stringArgFunc StringArgFunc
}

func NewStringFunctionExecutor(name string, stringArgFunc StringArgFunc) StringFunctionExecutor {
	return StringFunctionExecutor{name, stringArgFunc}
}

func (afs StringFunctionExecutor) Process(input string) (dto.OutputMessage, output.ExecutionOutput) {
	err := afs.stringArgFunc(input)
	if err != nil {
		return nil, output.NonRecoverableError{Err: err}
	}
	
	return dto.StringOutputMessage(""), output.NonOutputtableSuccess{}
}

func (afs StringFunctionExecutor) GetFuncName() string {
	return afs.name
}