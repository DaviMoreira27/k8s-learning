package main

import (
	"fmt"
	"log"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func loadDBConfig() DBConfig {
	cfg := DBConfig{
		Host:     requireEnv("DATABASE_HOST"),
		Port:     requireEnv("DATABASE_PORT"),
		Name:     requireEnv("DATABASE_NAME"),
		User:     requireEnv("DATABASE_USER"),
		Password: requireEnv("DATABASE_PASSWORD"),
	}
	return cfg
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Name, c.User, c.Password,
	)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %q is not set", key)
	}
	return v
}
