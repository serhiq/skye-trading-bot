package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

const TempPatch = "./assets/tmp/"
const PreviewCachePatch = "./assets/previews/"
const FileProviderPatch = "./assets/images/"

var (
	version string = "1.0"
)

type Config struct {
	Project    Project    `yaml:"project"`
	Telegram   Telegram   `yaml:"telegram"`
	ProductAPI ProductAPI `yaml:"product_api"`
	OrderAPI   OrderAPI   `yaml:"order_api"`
	DBConfig   DBConfig   `yaml:"database"`
	Telemetry  Telemetry  `yaml:"telemetry"`
}

type Project struct {
	Name        string `yaml:"name"`
	ServiceName string `yaml:"serviceName"`
	Version     string
	Timezone    string `yaml:"timezone"`
}

type Telegram struct {
	Token string `yaml:"token" envconfig:"TELEGRAM_TOKEN" validate:"required"`
}

const (
	EvotorAPIKind ProductApiKind = "evo_api"
	RestoAPIKind  ProductApiKind = "resto_api"
	FileKind      ProductApiKind = "file_provider"
)

type ProductApiKind string

type ProductAPI struct {
	Kind     ProductApiKind `yaml:"kind"   envconfig:"PRODUCT_API_KIND,omitempty"`
	BaseURL  string         `yaml:"base_url" envconfig:"PRODUCT_API_BASE_URL,omitempty"`
	Auth     string         `yaml:"auth" envconfig:"PRODUCT_API_AUTH,omitempty"`
	Store    string         `yaml:"store" envconfig:"PRODUCT_API_STORE,omitempty"`
	MenuUuid string         `yaml:"menu_root" envconfig:"PRODUCT_API_ROOT_UUID,omitempty"`
}

const (
	OrderAPIKind      OrderApiKind = "order_api"
	RestoOrderAPIKind OrderApiKind = "resto_api"
	FileOrderKind     OrderApiKind = "file_provider"
)

type OrderApiKind string

type OrderAPI struct {
	Kind    OrderApiKind `yaml:"kind"     envconfig:"ORDER_API_KIND,omitempty"`
	BaseURL string       `yaml:"base_url" envconfig:"ORDER_API_BASE_URL,omitempty"`
	Auth    string       `yaml:"auth"     envconfig:"ORDER_API_AUTH,omitempty"`
	Store   string       `yaml:"store"    envconfig:"ORDER_API_STORE,omitempty"`
}

type DBConfig struct {
	Host         string `yaml:"host" envconfig:"DB_HOST"`
	Port         int    `yaml:"port" envconfig:"DB_PORT"`
	DatabaseName string `yaml:"database_name" envconfig:"DB_DATABASE_NAME"`
	Username     string `yaml:"username" envconfig:"DB_USERNAME"`
	Password     string `yaml:"password" envconfig:"DB_PASSWORD"`
}

type Telemetry struct {
	GraylogPath string `yaml:"graylogPath"`
}

func New() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./configs/config.yaml"
	}

	config := &Config{}

	if err := fromYaml(path, config); err != nil {
		fmt.Printf("couldn'n load config from %s: %s\r\n", path, err.Error())
	}

	if err := fromEnv(config); err != nil {
		fmt.Printf("couldn'n load config from env: %s\r\n", err.Error())
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(TempPatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating tmp path %s (%s)", filepath.Dir(TempPatch), err)
	}

	if err := os.MkdirAll(filepath.Dir(PreviewCachePatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating cache path %s (%s)", filepath.Dir(PreviewCachePatch), err)
	}

	config.Project.Version = version

	return config, nil
}

func fromYaml(path string, config *Config) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func fromEnv(config *Config) error {
	return envconfig.Process("", config)
}

func validate(cfg *Config) error {
	if cfg.Telegram.Token == "" {
		return fmt.Errorf("config: %s is not set", "TELEGRAM_TOKEN")
	}

	if cfg.ProductAPI.Kind == "" {
		return fmt.Errorf("config: %s is not set", "Kind")
	}

	if cfg.ProductAPI.Kind != FileKind {
		if cfg.ProductAPI.BaseURL == "" {
			return fmt.Errorf("config: %s is not set", "API_BASE_URL")
		}

		if cfg.ProductAPI.Auth == "" {
			return fmt.Errorf("config: %s is not set", "Auth")
		}

		if cfg.ProductAPI.Store == "" {
			return fmt.Errorf("config: %s is not set", "Store")
		}
	}

	if cfg.OrderAPI.Kind == "" {
		return fmt.Errorf("config: %s is not set", "Kind")
	}

	if cfg.OrderAPI.Kind != FileOrderKind {
		if cfg.OrderAPI.BaseURL == "" {
			return fmt.Errorf("config: %s is not set", "API_BASE_URL")
		}

		if cfg.OrderAPI.Auth == "" {
			return fmt.Errorf("config: %s is not set", "Auth")
		}

		if cfg.OrderAPI.Store == "" {
			return fmt.Errorf("config: %s is not set", "Store")
		}
	}

	if cfg.DBConfig.DatabaseName == "" {
		return fmt.Errorf("config: %s is not set", "DB_DATABASE_NAME")
	}

	if cfg.DBConfig.Username == "" {
		return fmt.Errorf("config: %s is not set", "DB_USERNAME")
	}

	return nil
}
