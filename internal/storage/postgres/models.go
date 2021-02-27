package postgres

import "time"

type Organisation struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerID string
}

type Schedule struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerID        string
	OrganisationID string
}

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserExternalLogin struct {
	ID        string
	Source    string
	SourceID  string
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID string
}
