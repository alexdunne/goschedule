package accounts

type Schedule struct {
	ID   string
	Name string

	OrganisationID string
	OwnerID        string
}

func (o *Schedule) Validate() error {
	return nil
}
