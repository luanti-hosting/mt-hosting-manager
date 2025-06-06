package core

import (
	"context"
	"fmt"
	"io"
	"mt-hosting-manager/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("%s.zip", b.ID)
}

func (c *Core) GetS3Client() (*minio.Client, error) {
	return minio.New(c.cfg.S3_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(c.cfg.S3_KEYID, c.cfg.S3_ACCESSKEY, ""),
		Secure: true,
	})
}

func (c *Core) RemoveBackup(b *types.Backup) error {
	client, err := c.GetS3Client()
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}

	err = client.RemoveObject(context.Background(), c.cfg.S3_BUCKET, getBackupFilename(b), minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3 remove error: %v", err)
	}

	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) StreamBackup(b *types.Backup, w io.Writer) error {
	client, err := c.GetS3Client()
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}

	// TODO

	r, err := client.ReadStream(getBackupFilename(b))
	if err != nil {
		return fmt.Errorf("readstream error: %v", err)
	}
	defer r.Close()

	var reader io.Reader
	reader = r

	if b.Passphrase != "" {
		// enable decryption
		reader, err = EncryptedReader(b.Passphrase, reader)
		if err != nil {
			return fmt.Errorf("decryption failed: %v", err)
		}
	}

	_, err = io.Copy(w, reader)
	return err
}
