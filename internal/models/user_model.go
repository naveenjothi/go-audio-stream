package models

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
	Username  string `json:"user_name"`

	// add User follows and so on
}
