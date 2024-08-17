package responses


type UserResponseRegister struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

type UserResponseLogin struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Token string `json:"token"`
}