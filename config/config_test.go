package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test setup
	testCases := []struct {
		name       string
		configYaml string
		envVars    map[string]string
		expected   *Config
	}{
		{
			name: "basic config from file",
			configYaml: `
database:
  driver: postgres
  host: localhost
  port: "5432"
  username: user
  password: pass
  database: testdb
http_server:
  port: "8080"
  host: localhost
  ssl: false
  cert_file: cert.pem
  key_file: key.pem
smtp_server:
  host: smtp.example.com
  port: "587"
  username: smtp_user
  password: smtp_pass
`,
			expected: &Config{
				Database: Database{
					Driver:   "postgres",
					Host:     "localhost",
					Port:     "5432",
					Username: "user",
					Password: "pass",
					Database: "testdb",
				},
				HttpServer: HttpServer{
					Port:     "8080",
					Host:     "localhost",
					SSL:      false,
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
				SmtpServer: SmtpServer{
					Host:     "smtp.example.com",
					Port:     "587",
					Username: "smtp_user",
					Password: "smtp_pass",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary config file
			tmpfile, err := os.CreateTemp("", "config*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write the test config to the temporary file
			if _, err := tmpfile.Write([]byte(tc.configYaml)); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			// Set up environment variables if any
			for k, v := range tc.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Reset viper before each test
			viper.Reset()

			// Run the test
			config := Load(tmpfile.Name())

			// Assertions
			assert.Equal(t, tc.expected.Database.Driver, config.Database.Driver)
			assert.Equal(t, tc.expected.Database.Host, config.Database.Host)
			assert.Equal(t, tc.expected.Database.Port, config.Database.Port)
			assert.Equal(t, tc.expected.Database.Username, config.Database.Username)
			assert.Equal(t, tc.expected.Database.Password, config.Database.Password)
			assert.Equal(t, tc.expected.Database.Database, config.Database.Database)

			assert.Equal(t, tc.expected.HttpServer.Port, config.HttpServer.Port)
			assert.Equal(t, tc.expected.HttpServer.Host, config.HttpServer.Host)
			assert.Equal(t, tc.expected.HttpServer.SSL, config.HttpServer.SSL)
			assert.Equal(t, tc.expected.HttpServer.CertFile, config.HttpServer.CertFile)
			assert.Equal(t, tc.expected.HttpServer.KeyFile, config.HttpServer.KeyFile)

			assert.Equal(t, tc.expected.SmtpServer.Host, config.SmtpServer.Host)
			assert.Equal(t, tc.expected.SmtpServer.Port, config.SmtpServer.Port)
			assert.Equal(t, tc.expected.SmtpServer.Username, config.SmtpServer.Username)
			assert.Equal(t, tc.expected.SmtpServer.Password, config.SmtpServer.Password)
		})
	}
}

func TestLoadConfigError(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Test that Load() panics with non-existent config
	assert.Panics(t, func() {
		// Point to a non-existent config file
		Load("nonexistent.yaml")
	})
}
