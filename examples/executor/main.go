package main

import (
	"async_executor/examples/builder"
	"github.com/breathbath/go_utils/utils/errs"
)

func main() {
	asyncExecutor, err := builder.BuildAsyncFuncExecutionFacade()
	errs.FailOnError(err)

	err = asyncExecutor.ExecuteAsyncFunctions()
	errs.FailOnError(err)
}
