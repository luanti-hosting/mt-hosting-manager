package core

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/types"
)

func (c *Core) CheckHetznerServerTypes() error {
	server_types, err := c.hcloud.GetServerTypes()
	if err != nil {
		return fmt.Errorf("hcloud call error: %v", err)
	}

	server_type_map := map[string]*hetzner_cloud.ServerType{}
	for _, st := range server_types.ServerTypes {
		server_type_map[st.Name] = st
	}

	available_server_types, err := c.repos.NodeTypeRepo.GetAll()
	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}

	for _, st := range available_server_types {
		if st.Provider != types.ProviderHetzner {
			continue
		}

		hst := server_type_map[st.ServerType]
		if hst == nil {
			return fmt.Errorf("server not found in hetzner list: '%s'", st.ServerType)
		}

		if hst.Deprecated {
			return fmt.Errorf("Server-type deprecated: '%s'", st.ServerType)
		}
	}

	return nil
}
