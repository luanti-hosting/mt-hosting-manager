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
		return fmt.Errorf("create s3 client error: %v", err)
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
		return fmt.Errorf("create s3 client error: %v", err)
	}

	obj, err := client.GetObject(context.Background(), c.cfg.S3_BUCKET, getBackupFilename(b), minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("client getObject error: %v", err)
	}

	defer obj.Close()

	_, err = io.Copy(w, obj)
	return err
}

func (c *Core) StoreBackup(b *types.Backup, r io.Reader) (int64, error) {
	client, err := c.GetS3Client()
	if err != nil {
		return 0, fmt.Errorf("create s3 client error: %v", err)
	}

	info, err := client.PutObject(context.Background(), c.cfg.S3_BUCKET, getBackupFilename(b), r, -1, minio.PutObjectOptions{})
	if err != nil {
		return 0, fmt.Errorf("client putObject error: %v", err)
	}

	return info.Size, nil
}
