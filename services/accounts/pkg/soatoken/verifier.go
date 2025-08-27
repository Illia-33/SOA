package soatoken

type RightsRequirements struct {
	Read  bool
	Write bool
}

type Verifier interface {
	Verify(token string, r RightsRequirements) error
}
