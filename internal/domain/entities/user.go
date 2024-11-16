package entities

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"registeredAt"`
}
