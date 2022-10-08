package chat

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Id       int    `json:"id" db:"id"`
	Date     string `json:"date" db:"date"`
	Username string `json:"username" db:"username"`
	Body     string `json:"message" db:"message"`
}
