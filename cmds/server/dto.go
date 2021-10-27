package server

type userDto struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

type loginDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type tokenDto struct {
	Token string `json:"token"`
}
