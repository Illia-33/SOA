package api

type ErrorNoUserId struct{}
type ErrorTooMuchUserId struct{}
type ErrorNegativeTtl struct{}

func (ErrorNoUserId) Error() string {
	return "no userid"
}

func (ErrorTooMuchUserId) Error() string {
	return "too much userid"
}

func (ErrorNegativeTtl) Error() string {
	return "ttl is negative"
}
