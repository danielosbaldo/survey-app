package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "postgres")
	pass := env("DB_PASS", "postgres")
	name := env("DB_NAME", "heladeria")
	ssl  := env("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", host, user, pass, name, port, ssl)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
