package types

type ServiceTicketState string

const (
	ServiceTicketOpen     ServiceTicketState = "OPEN"
	ServiceTicketResolved ServiceTicketState = "RESOLVED"
	ServiceTicketClosed   ServiceTicketState = "CLOSED"
)

type ServiceTicket struct {
	ID               string             `json:"id" gorm:"primarykey;column:id"`
	Title            string             `json:"title" gorm:"column:title"`
	UserID           string             `json:"user_id" gorm:"column:user_id"`
	UserNodeID       *string            `json:"user_node_id" gorm:"column:user_node_id"`
	MinetestServerID *string            `json:"minetest_server_id" gorm:"column:minetest_server_id"`
	BackupID         *string            `json:"backup_id" gorm:"column:backup_id"`
	Created          int64              `json:"created" gorm:"column:created"`
	Closed           *int64             `json:"closed" gorm:"column:closed"`
	State            ServiceTicketState `json:"state" gorm:"column:state"`
}

func (*ServiceTicket) TableName() string {
	return "service_ticket"
}

type ServiceTicketSearch struct {
	TicketID         *string             `json:"ticket_id"`
	UserID           *string             `json:"user_id"`
	UserNodeID       *string             `json:"user_node_id"`
	MinetestServerID *string             `json:"minetest_server"`
	BackupID         *string             `json:"backup_id"`
	State            *ServiceTicketState `json:"state"`
}

type ServiceTicketMessage struct {
	ID        string `json:"id" gorm:"primarykey;column:id"`
	TicketID  string `json:"ticket_id" gorm:"column:ticket_id"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
	Message   string `json:"message" gorm:"column:message"`
}

func (*ServiceTicketMessage) TableName() string {
	return "service_ticket_message"
}
