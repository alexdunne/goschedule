package accounts

const (
	AuthSourceGitHub = "github"
)

type Account struct {
	Name     string
	Email    string
	Source   string
	SourceID string

	UserID string
}

func (a *Account) Validate() error {
	return nil
}
