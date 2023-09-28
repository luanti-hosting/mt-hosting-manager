package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeDiscord UserType = "DISCORD"
	UserTypeMesehub UserType = "MESEHUB"
)

type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type UserState string

const (
	UserStateActive UserState = "ACTIVE"
)

func UserProvider() *User { return &User{} }

type User struct {
	ID             string    `json:"id"`
	State          UserState `json:"state"`
	Name           string    `json:"name"`
	Hash           string    `json:"hash"`
	Mail           string    `json:"mail"`
	MailVerified   bool      `json:"mail_verified"`
	ActivationCode string    `json:"activation_code"`
	Created        int64     `json:"created"`
	Balance        int64     `json:"balance"`
	WarnBalance    int64     `json:"warn_balance"`
	ExternalID     string    `json:"external_id"`
	Type           UserType  `json:"type"`
	Role           UserRole  `json:"role"`
}

type UserSearch struct {
	MailLike *string `json:"mail_like"`
	Limit    *int    `json:"limit"`
}

func (m *User) Columns(action string) []string {
	return []string{"id", "state", "name", "hash", "mail", "mail_verified", "activation_code", "created", "balance", "warn_balance", "external_id", "type", "role"}
}

func (m *User) Table() string {
	return "user"
}

func (m *User) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.State, &m.Name, &m.Hash, &m.Mail, &m.MailVerified, &m.ActivationCode, &m.Created, &m.Balance, &m.WarnBalance, &m.ExternalID, &m.Type, &m.Role)
}

func (m *User) Values(action string) []any {
	return []any{m.ID, m.State, m.Name, m.Hash, m.Mail, m.MailVerified, m.ActivationCode, m.Created, m.Balance, m.WarnBalance, m.ExternalID, m.Type, m.Role}
}
