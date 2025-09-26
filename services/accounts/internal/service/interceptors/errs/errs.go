package errs

type UnknownAuthKind struct{}
type NoAuth struct{}
type InvalidToken struct{}
type NoMetadata struct{}

func (UnknownAuthKind) Error() string {
	return "unknown auth kind"
}
func (NoAuth) Error() string {
	return "auth required"
}

func (InvalidToken) Error() string {
	return "token is invalid"
}

func (NoMetadata) Error() string {
	return "metadata not found"
}
