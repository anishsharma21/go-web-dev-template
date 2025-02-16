package models

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Email     string    `json:"email"`
	Password  *string   `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
