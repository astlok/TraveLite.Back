package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	LogLevel string `yaml:"log_level"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Workdir  string `yaml:"workdir"`
	Name     string `yaml:"name"`

	Postgres struct {
		User     string `yaml:"user"`
		DBName   string `yaml:"db_name"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"postgres"`

	Cors struct {
		AllowMethods     []string      `yaml:"allow_methods"`
		AllowCredentials bool          `yaml:"allow_credentials"`
		AllowHeaders     []string      `yaml:"allow_headers"`
		AllowFiles       bool          `yaml:"allow_files"`
		AllowOrigins     []string      `yaml:"allow_origins"`
		MaxAge           time.Duration `yaml:"max_age"`
	} `yaml:"cors"`

	S3 struct {
		Bucket  string `yaml:"bucket"`
		FileURL string `yaml:"file_url"`
	} `yaml:"s3"`
}

func ParseConfig(filePath string) (*Config, error) {
	var config Config

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(nil, "can't read from file")
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal config to yaml")
	}

	return &config, nil
}
