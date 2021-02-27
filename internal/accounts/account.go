package accounts

const (
	AuthSourceGitHub = "github"
)

type NewAccount struct {
	Name     string
	Email    string
	Source   string
	SourceID string

	UserID string
}

func (a *NewAccount) Validate() error {
	return nil
}
