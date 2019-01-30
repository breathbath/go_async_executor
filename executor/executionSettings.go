package executor

//ExecutionSettings set of options for ExecutionFacade
type ExecutionSettings struct {
	ExecutionProcessorsCount int //how many routines to run for executing async payloads
	PublishBadlyFormattedMessagesToErrorChannel bool //shall it send non-convertable messages to the error outputter
	FailedMessagesRepeatAttemptsCount int //how many times to repeat failed executions
}
