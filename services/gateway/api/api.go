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
