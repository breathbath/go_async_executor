package main

import (
	"go_async_executor/examples/builder"
	"github.com/breathbath/go_utils/utils/errs"
)

func main() {
	asyncFuncRegistrationFacade, err := builder.BuildAsyncFuncRegistrationFacade()
	errs.FailOnError(err)

	err = asyncFuncRegistrationFacade.CallFunctionAsync("time_executor", "Call me async", 0)
	errs.FailOnError(err)

	err = asyncFuncRegistrationFacade.CallFunctionAsync("fail_me", "", 0)
	errs.FailOnError(err)
}
