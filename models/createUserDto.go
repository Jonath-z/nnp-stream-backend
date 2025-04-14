package models

type CreateUserDto struct {
	Email       string `json:"email"`
	ClerkUserId string `json:"clerk_id"`
}
