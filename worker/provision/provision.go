package provision

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ProvisionModel struct {
	Config   *types.Config
	UserNode *types.UserNode
}

func Provision(client *ssh.Client, cfg *types.Config, node *types.UserNode) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("could not open session: %v", err)
	}
	defer session.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("could not create sftp client: %v", err)
	}
	defer sftp.Close()

	dirs := []string{
		"/etc/docker",
		"/etc/iptables",
		"/provision",
	}
	for _, dir := range dirs {
		err = core.SCPMkDir(sftp, dir)
		if err != nil {
			return err
		}
	}

	err = core.SCPTemplateFile(sftp, Files, "daemon.json", "/etc/docker/daemon.json", 0644, true, nil)
	if err != nil {
		return err
	}

	err = core.SCPTemplateFile(sftp, Files, "rules.v6", "/etc/iptables/rules.v6", 0644, true, nil)
	if err != nil {
		return err
	}

	err = core.SCPWriteFile(sftp, Files, "setup.sh", "/provision/setup.sh", 0755, true)
	if err != nil {
		return fmt.Errorf("could not write setup.sh: %v", err)
	}

	m := &ProvisionModel{
		Config:   cfg,
		UserNode: node,
	}

	err = core.SCPTemplateFile(sftp, Files, "docker-compose.yml", "/provision/docker-compose.yml", 0644, true, m)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	_, stderr, err := core.SSHExecute(client, "/provision/setup.sh")
	if err != nil {
		return fmt.Errorf("SSHExecute error: %v, stderr: '%s'", err, string(stderr))
	}

	return nil
}
