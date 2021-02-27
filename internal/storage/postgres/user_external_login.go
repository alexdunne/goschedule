package postgres

import "time"

type UserExternalLogin struct {
	ID        string
	Source    string
	SourceID  string
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID string
}
