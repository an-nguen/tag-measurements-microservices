package dto

// Login credentials for remote api service
type LoginDataRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
