package errs

type AccessDenied struct{}
type NoWriteAccess struct{}
type NoReadAccess struct{}
type TokenExpired struct{}

func (AccessDenied) Error() string {
	return "access denied"
}

func (NoWriteAccess) Error() string {
	return "no write access"
}

func (NoReadAccess) Error() string {
	return "no read access"
}

func (TokenExpired) Error() string {
	return "token expired"
}
