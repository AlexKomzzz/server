package chat

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RowChat struct {
	Id       int
	Data     string
	Username string
	Message  string
}
