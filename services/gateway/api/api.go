package api

type RegisterProfileRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
}

type RegisterProfileResponse struct {
	ProfileID string `json:"profile_id"`
}

type GetProfileResponse struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
}

type EditProfileRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Birthday    string `json:"birthday"`
	Bio         string `json:"bio"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type AuthenticateRequest struct {
	Login       string `json:"login"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
}

type CreateApiTokenRequest struct {
	Auth        AuthenticateRequest `json:"auth"`
	ReadAccess  bool                `json:"read_access"`
	WriteAccess bool                `json:"write_access"`
	Ttl         string              `json:"ttl"`
}

type CreateApiTokenResponse struct {
	Token string `json:"token"`
}
