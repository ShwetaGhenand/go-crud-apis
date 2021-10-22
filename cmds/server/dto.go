package server

type userDto struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}
