package errs

type ProfileNotFound struct{}
type TokenNotFound struct{}
type PasswordsDoNotMatch struct{}
type AccountNotFound struct{}
type UserIdNotFound struct{}

func (ProfileNotFound) Error() string {
	return "profile not found"
}

func (TokenNotFound) Error() string {
	return "token not found"
}

func (PasswordsDoNotMatch) Error() string {
	return "passwords do not match"
}

func (AccountNotFound) Error() string {
	return "account not found"
}

func (UserIdNotFound) Error() string {
	return "user id not found"
}
