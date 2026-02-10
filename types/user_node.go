package types

type UserNodeState string

const (
	UserNodeStateCreated        UserNodeState = "CREATED"
	UserNodeStateProvisioning   UserNodeState = "PROVISIONING"
	UserNodeStateRunning        UserNodeState = "RUNNING"
	UserNodeStateRemoving       UserNodeState = "REMOVING"
	UserNodeStateDecommissioned UserNodeState = "DECOMMISSIONED"
)

// Created -> Provisioning -> Running -> Removing -> Decommissioned

type UserNode struct {
	ID                  string        `json:"id" gorm:"primarykey;column:id"`
	UserID              string        `json:"user_id" gorm:"column:user_id"`
	NodeTypeID          string        `json:"node_type_id" gorm:"column:node_type_id"`
	Location            string        `json:"location" gorm:"column:location"`
	ExternalID          string        `json:"external_id" gorm:"column:external_id"`
	Created             int64         `json:"created" gorm:"column:created"`
	ValidUntil          int64         `json:"valid_until" gorm:"column:valid_until"`
	State               UserNodeState `json:"state" gorm:"column:state"`
	Name                string        `json:"name" gorm:"column:name"`
	Alias               string        `json:"alias" gorm:"column:alias"`
	IPv4                string        `json:"ipv4" gorm:"column:ipv4"`
	IPv6                string        `json:"ipv6" gorm:"column:ipv6"`
	ExternalIPv4DNSID   string        `json:"external_ipv4_dns_id" gorm:"column:external_ipv4_dns_id"`
	ExternalIPv6DNSID   string        `json:"external_ipv6_dns_id" gorm:"column:external_ipv6_dns_id"`
	NodeExporterVersion string        `json:"node_exporter_version" gorm:"column:node_exporter_version"`
	TraefikVersion      string        `json:"traefik_version" gorm:"column:traefik_version"`
	IPv6NatVersion      string        `json:"ipv6nat_version" gorm:"column:ipv6nat_version"`
	Fingerprint         string        `json:"fingerprint" gorm:"column:fingerprint"`
}

func (m *UserNode) TableName() string {
	return "user_node"
}

type UserNodeSearch struct {
	ID         *string        `json:"id"`
	Name       *string        `json:"name"`
	UserID     *string        `json:"user_id"`
	State      *UserNodeState `json:"state"`
	ValidUntil *int64         `json:"valid_until"`
}
