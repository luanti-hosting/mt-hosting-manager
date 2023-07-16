package types

type UserNodeState string

const (
	UserNodeStateCreated UserNodeState = "CREATED"
)

type UserNode struct {
	ID         string        `json:"id"`
	UserID     string        `json:"user_id"`
	NodeTypeID string        `json:"node_type_id"`
	Created    int64         `json:"created"`
	State      UserNodeState `json:"state"`
	Name       string        `json:"name"`
	IPv4       string        `json:"ipv4"`
	IPv6       string        `json:"ipv6"`
}

func (m *UserNode) Columns(action string) []string {
	return []string{
		"id",
		"user_id",
		"node_type_id",
		"created",
		"state",
		"name",
		"ipv4",
		"ipv6",
	}
}

func (m *UserNode) Table() string {
	return "user_node"
}

func (m *UserNode) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.UserID, &m.NodeTypeID, &m.Created, &m.State, &m.Name, &m.IPv4, &m.IPv6)
}

func (m *UserNode) Values(action string) []any {
	return []any{m.ID, m.UserID, m.NodeTypeID, m.Created, m.State, m.Name, m.IPv4, m.IPv6}
}
