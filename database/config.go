package database

import "fmt"

// Config represents database configuration.
type Config struct {
	Host             string `mapstructure:"host"`       // server address
	Port             int    `mapstructure:"port"`       // server port
	Username         string `mapstructure:"username"`   // user
	Password         string `mapstructure:"password"`   // pass
	Database         string `mapstructure:"database"`   // database
	MigrationsSource string `mapstructure:"migrations"` // database
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.Database)
}
