package requests

type RegisterRequest struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
