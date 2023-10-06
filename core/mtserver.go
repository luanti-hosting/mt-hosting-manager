package core

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/types"
	"time"
)

type CreateServerResult struct {
	Valid                 bool `json:"valid"`
	PortInvalid           bool `json:"port_invalid"`
	PortUsed              bool `json:"port_used"`
	AdminNameInvalid      bool `json:"admin_name_invalid"`
	ServerNameInvalid     bool `json:"server_name_invalid"`
	ServerNameAlreadyUsed bool `json:"server_name_used"`
}

var hdns_records *hetzner_dns.RecordsResponse
var hdns_records_updated time.Time

func (c *Core) ValidateCreateServer(server *types.MinetestServer, node *types.UserNode) (*CreateServerResult, error) {
	csr := &CreateServerResult{
		Valid: true,
	}

	other_servers, err := c.repos.MinetestServerRepo.GetByNodeID(node.ID)
	if err != nil {
		return nil, fmt.Errorf("servers fetch error: %v", err)
	}

	if server.Port < 1000 || server.Port > 65000 {
		csr.PortInvalid = true
		csr.Valid = false
	}

	err = types.ValidateUsername(server.Admin)
	if err != nil {
		csr.AdminNameInvalid = true
		csr.Valid = false
	}

	for _, s := range other_servers {
		if s.State == types.MinetestServerStateDecommissioned {
			continue
		}
		if s.Port == server.Port {
			csr.PortUsed = true
			csr.Valid = false
			break
		}
	}

	if !types.ValidServerName.Match([]byte(server.DNSName)) {
		csr.ServerNameInvalid = true
		csr.Valid = false
	}

	if time.Since(hdns_records_updated) > 5*time.Minute {
		// fetch records
		hdns_records, err = c.hdns.GetRecords()
		if err != nil {
			return nil, fmt.Errorf("error in hetzner dns api: %v", err)
		}
		hdns_records_updated = time.Now()
	}

	for _, existing_record := range hdns_records.Records {
		if existing_record.Name == server.DNSName {
			csr.ServerNameAlreadyUsed = true
			csr.Valid = false
			break
		}
	}

	return csr, nil
}
