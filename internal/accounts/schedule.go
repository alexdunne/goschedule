package accounts

type Schedule struct {
	ID      string
	Name    string
	OwnerID string
}

func (o *Schedule) Validate() error {
	return nil
}
