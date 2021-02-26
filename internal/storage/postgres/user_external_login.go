package postgres

import "time"

type UserExternalLogin struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Source    string    `json:"source"`
	SourceID  string    `json:"sourceId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
