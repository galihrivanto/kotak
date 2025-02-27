package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type configKey struct{}

var contextKey = configKey{}

// Logger configuration
type Logger struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
}

// Database configuration
type Database struct {
	Driver   string `mapstructure:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Database string `mapstructure:"database" yaml:"database"`
}

func (d Database) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.Username, d.Password, d.Database)
}

// HttpServer configuration
type HttpServer struct {
	Port     string `mapstructure:"port" yaml:"port"`
	Host     string `mapstructure:"host" yaml:"host"`
	TLS      bool   `mapstructure:"tls" yaml:"tls"`
	CertFile string `mapstructure:"cert_file" yaml:"cert_file"`
	KeyFile  string `mapstructure:"key_file" yaml:"key_file"`

	// API configuration
	APIHost string `mapstructure:"api_host" yaml:"api_host"`
	APIBase string `mapstructure:"api_base" yaml:"api_base"`

	// Static file configuration
	StaticDir string `mapstructure:"static_dir" yaml:"static_dir"`
	StaticURL string `mapstructure:"static_url" yaml:"static_url"`
}

// SmtpServer configuration
type SmtpServer struct {
	AppName  string `mapstructure:"app_name" yaml:"app_name"`
	Hostname string `mapstructure:"hostname" yaml:"hostname"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	MaxSize  int    `mapstructure:"max_size" yaml:"max_size"`
}

// Config is the configuration for the application
type Config struct {
	Database   Database   `mapstructure:"database" yaml:"database"`
	HttpServer HttpServer `mapstructure:"http_server" yaml:"http_server"`
	SmtpServer SmtpServer `mapstructure:"smtp_server" yaml:"smtp_server"`
	Logger     Logger     `mapstructure:"logger" yaml:"logger"`
}

// Load loads the configuration from the given viper instance
func Load(vars ...string) *Config {
	viper.SetConfigType("yaml")
	if len(vars) == 0 {
		fmt.Println("Loading config from default file")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	} else {
		fmt.Println("Loading config from file:", vars[0])
		viper.SetConfigFile(vars[0])
	}

	viper.AutomaticEnv()

	// Replace dots with underscores in env variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return config
}

// FromContext returns the Config instance from the context
func FromContext(ctx context.Context) *Config {
	return ctx.Value(contextKey).(*Config)
}

// WithContext returns a new context with the Config instance
func WithContext(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, contextKey, config)
}
