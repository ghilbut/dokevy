package config

import (
	"os"

	// external
	log "github.com/sirupsen/logrus"
)

type StorageType string

const (
	StorageTypeGCS StorageType = "GCS"
	StorageTypeS3  StorageType = "S3"
)

type Config struct {
	LogLevel    log.Level
	Address     string
	Bucket      string
	StorageType StorageType
	GCS         *GCSConfig
	S3          *S3Config
}

func NewConfig() *Config {
	config := &Config{
		LogLevel:    log.InfoLevel,
		Address:     os.Getenv("LISTEN_ADDRESS"),
		StorageType: StorageTypeS3,
	}

	switch config.StorageType {
	case StorageTypeGCS:
		config.GCS = &GCSConfig{}
	case StorageTypeS3:
		config.Bucket = os.Getenv("AWS_S3_BUCKET")
		config.S3 = &S3Config{
			AccessKey: os.Getenv("AWS_ACCESS_KEY"),
			SecretKey: os.Getenv("AWS_SECRET_KEY"),
			Token:     os.Getenv("AWS_TOKEN"),
			Region:    os.Getenv("AWS_S3_REGION"),
		}
	default:
		return nil
	}

	return config
}
