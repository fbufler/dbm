package database

import "context"

type Config struct {
	FilePath          string `mapstructure:"file_path"`
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`
	UseSSL            bool   `mapstructure:"use_ssl"`
	SSLCertPath       string `mapstructure:"ssl_cert_path"`
	SSLKeyPath        string `mapstructure:"ssl_key_path"`
	SSLRootCertPath   string `mapstructure:"ssl_root_cert_path"`
	ConnectionTimeout int    `mapstructure:"connection_timeout"`
}

type Database interface {
	Connect() error
	Close() error
	Identifier() string
	Test(ctx context.Context) error
	TestWrite(ctx context.Context) error
	TestRead(ctx context.Context) error
	SetupTestTable(ctx context.Context) error
}
