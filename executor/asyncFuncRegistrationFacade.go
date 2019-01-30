package executor

import (
	"async_executor/dto"
	"fmt"
	"github.com/breathbath/go_utils/utils/enc"
	"time"
)

//AsyncFuncRegistrationFacade main entry point to register async function calls
type AsyncFuncRegistrationFacade struct {
	funcExecutorRegistry *AsyncFuncExecutorRegistry
	asyncFuncRegistrator AsyncFuncRegistrator
}

func NewAsyncFuncRegistrationFacade(
	funcExecutorRegistry *AsyncFuncExecutorRegistry,
	asyncFuncRegistrator AsyncFuncRegistrator,
) *AsyncFuncRegistrationFacade {
	return &AsyncFuncRegistrationFacade{
		funcExecutorRegistry: funcExecutorRegistry,
		asyncFuncRegistrator: asyncFuncRegistrator,
	}
}

func (cf *AsyncFuncRegistrationFacade) CallFunctionAsync(funcName, payload string, lifeTime time.Duration) error {
	if !cf.funcExecutorRegistry.HasFuncExecutor(funcName) {
		return fmt.Errorf("Unregistered async func %s", funcName)
	}

	input := dto.AsyncFuncInput{
		FunctionName:        funcName,
		Payload:             payload,
		CallsCount:          0,
		TimeStamp:           time.Now().UTC(),
		ValidWindow:         int64(lifeTime.Seconds()),
		MessageId:           enc.NewUuid(""),
		FailedAttemptsCount: 0,
		LastError:           "",
	}

	return cf.asyncFuncRegistrator.RegisterAsyncExecution(input)
}