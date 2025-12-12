package model

type User struct {
	ID         int    `json:"id" db:"id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"-" db:"password"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
}
