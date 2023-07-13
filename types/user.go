package types

type UserType string

const (
	UserTypeGithub  UserType = "GITHUB"
	UserTypeDiscord UserType = "DISCORD"
)

type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Mail        string   `json:"mail"`
	Created     int64    `json:"created"`
	ExternalID  string   `json:"external_id"`
	Type        UserType `json:"type"`
	Role        UserRole `json:"role"`
	Credits     int64    `json:"credits"`
	MaxCredits  int64    `json:"max_credits"`
	WarnCredits int64    `json:"warn_credits"`
}

func (m *User) Columns(action string) []string {
	return []string{"id", "name", "mail", "created", "external_id", "type", "role", "credits", "max_credits", "warn_credits"}
}

func (m *User) Table() string {
	return "user"
}

func (m *User) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Name, &m.Mail, &m.Created, &m.ExternalID, &m.Type, &m.Role, &m.Credits, &m.MaxCredits, &m.WarnCredits)
}

func (m *User) Values(action string) []any {
	return []any{m.ID, m.Name, m.Mail, m.Created, m.ExternalID, m.Type, m.Role, m.Credits, m.MaxCredits, m.WarnCredits}
}
