package core

import (
	"context"
	"fmt"
	"mt-hosting-manager/types"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type CreateServerResult struct {
	Valid                 bool `json:"valid"`
	PortInvalid           bool `json:"port_invalid"`
	PortUsed              bool `json:"port_used"`
	AdminNameInvalid      bool `json:"admin_name_invalid"`
	ServerNameInvalid     bool `json:"server_name_invalid"`
	ServerNameAlreadyUsed bool `json:"server_name_used"`
	ServerNameTooShort    bool `json:"server_name_too_short"`
	ServerNameReserved    bool `json:"server_name_reserved"`
}

var hdns_records []*hcloud.ZoneRRSet
var hdns_records_updated time.Time

func (c *Core) ValidateCreateServer(server *types.MinetestServer, node *types.UserNode) (*CreateServerResult, error) {
	csr := &CreateServerResult{
		Valid: true,
	}

	other_servers, err := c.repos.MinetestServerRepo.Search(&types.MinetestServerSearch{NodeID: &node.ID})
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
		zone, _, err := c.hc.Zone.GetByName(context.Background(), c.cfg.HetznerZoneName)
		if err != nil {
			return nil, fmt.Errorf("could not get zone '%s': %v", c.cfg.HetznerZoneName, err)
		}

		// fetch records
		hdns_records, _, err = c.hc.Zone.ListRRSets(context.Background(), zone, hcloud.ZoneRRSetListOpts{})
		if err != nil {
			return nil, fmt.Errorf("error in hetzner dns api: %v", err)
		}
		hdns_records_updated = time.Now()
	}

	// name of the record inside the zone (with suffix appended)
	record_name := fmt.Sprintf("%s%s", server.DNSName, c.cfg.DNSRecordSuffix)

	// check if the name is already used within the zone
	for _, existing_record := range hdns_records {
		if existing_record.Name == record_name {
			csr.ServerNameAlreadyUsed = true
			csr.Valid = false
			break
		}
	}

	// min length of 5 characters
	if len(server.DNSName) < 5 {
		csr.ServerNameTooShort = true
		csr.Valid = false
	}

	for _, resname := range c.cfg.ReservedPrefixes {
		if resname == server.DNSName {
			csr.ServerNameReserved = true
			csr.Valid = false
			break
		}
	}

	return csr, nil
}
