package core

import (
	"bytes"
	"embed"
	"fmt"
	"mt-hosting-manager/types"
	"net"
	"os"
	"text/template"
	"time"

	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func SCPWriteBytes(sftp *sftp.Client, data []byte, target string, mode os.FileMode, overwrite bool) error {
	if !overwrite {
		fi, _ := sftp.Stat(target)
		if fi != nil {
			// already exists
			return nil
		}
	}

	dstFile, err := sftp.Create(target)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	err = sftp.Chmod(target, mode)
	if err != nil {
		return fmt.Errorf("could not chmod file: %v", err)
	}

	return nil
}

func TemplateFile(fs embed.FS, filename string, model any) ([]byte, error) {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %v", filename, err)
	}

	t, err := template.New("tmpl").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("error templating %s: %v", filename, err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, model)
	if err != nil {
		return nil, fmt.Errorf("error executing template '%s': %v", filename, err)
	}

	return buf.Bytes(), nil
}

func SCPTemplateFile(sftp *sftp.Client, fs embed.FS, filename, target string, mode os.FileMode, overwrite bool, model any) error {
	data, err := TemplateFile(fs, filename, model)
	if err != nil {
		return err
	}

	err = SCPWriteBytes(sftp, data, target, mode, overwrite)
	if err != nil {
		return fmt.Errorf("template error in file '%s': %v", target, err)
	}

	return nil
}

func SCPWriteFile(sftp *sftp.Client, fs embed.FS, filename, target string, mode os.FileMode, overwrite bool) error {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return SCPWriteBytes(sftp, data, target, mode, overwrite)
}

func SCPMkDir(sftp *sftp.Client, dir string) error {
	_, err := sftp.Stat(dir)
	if err != nil {
		err = sftp.Mkdir(dir)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}
	return nil
}

func SSHExecute(client *ssh.Client, cmd string) ([]byte, []byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("could not open session: %v", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Start(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("start failed: %v", err)
	}

	err = session.Wait()
	if err != nil {
		ex, ok := err.(*ssh.ExitError)
		if ok {
			fmt.Printf("Exit status: %d\n", ex.ExitStatus())
			if ex.ExitStatus() != 0 {
				return stdout.Bytes(), stderr.Bytes(), fmt.Errorf("exit-status: %d", ex.ExitStatus())
			}
		} else {
			return stdout.Bytes(), stderr.Bytes(), fmt.Errorf("unknown script execution error: %v", err)
		}
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}

func CreateClient(node *types.UserNode) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:22", node.IPv4)
	key_file := os.Getenv("SSH_KEY")
	f, err := os.ReadFile(key_file)
	if err != nil {
		return nil, fmt.Errorf("ssh-key not found: %s", key_file)
	}

	key, err := ssh.ParsePrivateKey(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fp := ssh.FingerprintSHA256(key)
		if node.Fingerprint == "" {
			// no fingerprint yet, add and allow
			node.Fingerprint = fp
		}
		if fp != node.Fingerprint {
			// fingerprint mismatch
			return fmt.Errorf("fingerprint mismatch, on record: '%s' got: '%s'", node.Fingerprint, fp)
		}

		return nil
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}
	return client, nil
}

func TrySSHConnection(node *types.UserNode) (*ssh.Client, error) {
	try_count := 0
	for {
		client, err := CreateClient(node)
		if err != nil {
			if try_count > 5 {
				return nil, fmt.Errorf("ssh-client connection failed: %v", err)
			} else {
				logrus.WithFields(logrus.Fields{
					"err":       err,
					"try_count": try_count,
				}).Warn("ssh-client failed")
				try_count++
				time.Sleep(10 * time.Second)
			}
		} else {
			return client, nil
		}
	}
}
