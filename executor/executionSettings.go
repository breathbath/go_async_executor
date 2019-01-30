package executor

type ExecutionSettings struct {
	ExecutionProcessorsCount int
	PublishBadlyFormattedMessagesToErrorChannel bool
	FailedMessagesRepeatAttemptsCount int
}
