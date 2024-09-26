package entities_user_v1

import "time"

type User struct {
	ExternalID string    `json:"external_id"`
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Public_User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}
