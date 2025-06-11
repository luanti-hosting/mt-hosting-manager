package core

import (
	"context"
	"fmt"
	"io"
	"mt-hosting-manager/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (c *Core) getS3Client() (*minio.Client, error) {
	return minio.New(c.cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.cfg.S3KeyID, c.cfg.S3AccessKey, ""),
		Secure: true,
	})
}

func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("%s.zip", b.ID)
}

func (c *Core) RemoveBackup(b *types.Backup) error {
	client, err := c.getS3Client()
	if err != nil {
		return fmt.Errorf("create s3 client error: %v", err)
	}

	err = client.RemoveObject(context.Background(), c.cfg.S3Bucket, getBackupFilename(b), minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3 remove error: %v", err)
	}

	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) GetBackupSize(b *types.Backup) (int64, error) {
	client, err := c.getS3Client()
	if err != nil {
		return 0, fmt.Errorf("create s3 client error: %v", err)
	}

	fi, err := client.StatObject(context.Background(), c.cfg.S3Bucket, getBackupFilename(b), minio.GetObjectOptions{})
	if err != nil {
		return 0, fmt.Errorf("s3 stat error: %v", err)
	}

	return fi.Size, nil
}

func (c *Core) StreamBackup(b *types.Backup, w io.Writer) error {
	client, err := c.getS3Client()
	if err != nil {
		return fmt.Errorf("create s3 client error: %v", err)
	}

	r, err := client.GetObject(context.Background(), c.cfg.S3Bucket, getBackupFilename(b), minio.GetObjectOptions{})
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
