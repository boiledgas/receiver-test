package errors

type DisconnectedError struct {
}

type ParseError struct {
	message string
}

func (e DisconnectedError) Error() string {
	return "disconnected"
}

func (e ParseError) Error() string {
	return e.message
}

func NewDisconnectedError() error {
	return DisconnectedError{}
}

func New(message string) error {
	return ParseError{message: message}
}
