package storage

import (
	"crypto/tls"
	"go-starter-app/config"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinio(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Minio().Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio().AccessKey, cfg.Minio().SecretKey, ""),
		Secure: cfg.Minio().Secure,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	return minioClient, err
}
