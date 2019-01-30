package output

//NonRecoverableError if it is returned, message execution won't be repeated, message goes directly to error channel
type NonRecoverableError struct {
	Err error
}

func (nre NonRecoverableError) GetError() error {
	return nre.Err
}

//RecoverableError if it is returned, message execution will be be repeated,
//configurable amount of times after elapsing, message goes directly to error channel
// message goes directly to error channel
type RecoverableError struct {
	Err error
}

func (re RecoverableError) GetError() error {
	return re.Err
}
