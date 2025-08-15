package model

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
}
