package executor

//AsyncFuncExecutorRegistry contains functions to be called in async mode
type AsyncFuncExecutorRegistry struct {
	asyncFuncExecutors map[string]AsyncFunctionExecutor
}

func NewAsyncFuncExecutorRegistry() *AsyncFuncExecutorRegistry {
	return &AsyncFuncExecutorRegistry{
		asyncFuncExecutors: make(map[string]AsyncFunctionExecutor),
	}
}

func (afer *AsyncFuncExecutorRegistry) AddAsyncFunction(ip AsyncFunctionExecutor) {
	afer.asyncFuncExecutors[ip.GetFuncName()] = ip
}

func (afer *AsyncFuncExecutorRegistry) SetAsyncFunctions(processors map[string]AsyncFunctionExecutor) {
	afer.asyncFuncExecutors = processors
}

func (afer *AsyncFuncExecutorRegistry) HasFuncExecutor(funcName string) bool {
	_, ok := afer.asyncFuncExecutors[funcName]
	return ok
}

func (afer *AsyncFuncExecutorRegistry) GetFuncExecutor(funcName string) (AsyncFunctionExecutor, bool) {
	e, ok := afer.asyncFuncExecutors[funcName]

	return e, ok
}