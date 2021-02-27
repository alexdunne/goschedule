package accounts

type Organisation struct {
	ID      string
	Name    string
	OwnerID string
}

func (o *Organisation) Validate() error {
	return nil
}
