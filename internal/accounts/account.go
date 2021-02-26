package accounts

const (
	AuthSourceGitHub = "github"
)

type NewAccount struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Source   string `json:"source"`
	SourceID string `json:"sourceId"`

	UserID string
}
