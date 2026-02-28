package types

import (
	"fmt"
	"regexp"
)

type MinetestServerState string

const (
	MinetestServerStateCreated        MinetestServerState = "CREATED"
	MinetestServerStateProvisioning   MinetestServerState = "PROVISIONING"
	MinetestServerStateRunning        MinetestServerState = "RUNNING"
	MinetestServerStateRemoving       MinetestServerState = "REMOVING"
	MinetestServerStateDecommissioned MinetestServerState = "DECOMMISSIONED"
)

var ValidServerName = regexp.MustCompile(`^[a-z|A-Z|0-9]+$`)

var ValidPlayernameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`)

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("playername empty")
	}
	if len(username) > 20 {
		return fmt.Errorf("playername too long")
	}
	if !ValidPlayernameRegex.Match([]byte(username)) {
		return fmt.Errorf("playername can only contain chars a to z, A to Z, 0 to 9 and -, _")
	}
	return nil
}

type MinetestServer struct {
	ID                 string              `json:"id" gorm:"primarykey;column:id"`
	UserNodeID         string              `json:"user_node_id" gorm:"column:user_node_id"`
	Name               string              `json:"name" gorm:"column:name"`
	DNSName            string              `json:"dns_name" gorm:"column:dns_name"`
	Admin              string              `json:"admin" gorm:"column:admin"`
	ExternalCNAMEDNSID string              `json:"external_cname_dns_id" gorm:"column:external_cname_dns_id"`
	CustomDNS          string              `json:"custom_dns_name" gorm:"column:custom_dns_name"`
	Port               int                 `json:"port" gorm:"column:port"`
	UIVersion          string              `json:"ui_version" gorm:"column:ui_version"`
	NginxVersion       string              `json:"nginx_version" gorm:"column:nginx_version"`
	JWTKey             string              `json:"jwt_key" gorm:"column:jwt_key"`
	Created            int64               `json:"created" gorm:"column:created"`
	State              MinetestServerState `json:"state" gorm:"column:state"`
}

func (m *MinetestServer) TableName() string {
	return "minetest_server"
}

type MinetestServerSearch struct {
	ID     *string              `json:"id"`
	UserID *string              `json:"user_id"`
	NodeID *string              `json:"node_id"`
	State  *MinetestServerState `json:"state"`
}
