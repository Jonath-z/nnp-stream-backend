package models

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	ClerkId   string `json:"clerk_id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
