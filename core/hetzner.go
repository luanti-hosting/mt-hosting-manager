package core

import (
	"context"
	"fmt"
	"mt-hosting-manager/types"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func (c *Core) CheckHetznerServerTypes() error {
	server_types, err := c.hc.ServerType.All(context.Background())
	if err != nil {
		return fmt.Errorf("hcloud call error: %v", err)
	}

	server_type_map := map[string]*hcloud.ServerType{}
	for _, st := range server_types {
		server_type_map[st.Name] = st
	}

	available_server_types, err := c.repos.NodeTypeRepo.GetAll()
	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}

	for _, st := range available_server_types {
		if st.Provider != types.ProviderHetzner {
			continue // not a hetzner type
		}

		if st.State != types.NodeTypeStateActive {
			continue // not active
		}

		hst := server_type_map[st.ServerType]
		if hst == nil {
			return fmt.Errorf("server not found in hetzner list: '%s'", st.ServerType)
		}

		if hst.IsDeprecated() {
			return fmt.Errorf("Server-type deprecated: '%s'", st.ServerType)
		}
	}

	return nil
}
