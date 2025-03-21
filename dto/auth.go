package dto

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
