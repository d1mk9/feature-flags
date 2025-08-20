package config

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost string
	PostgresPort int
	PostgresDB   string

	PostgresUser     string
	PostgresPassword string
}

func (c Config) PostgresDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:   c.PostgresHost + ":" + strconv.Itoa(c.PostgresPort),
		Path:   "/" + c.PostgresDB,
	}
	q := url.Values{"sslmode": []string{"disable"}}
	u.RawQuery = q.Encode()
	return u.String()
}

func MustLoad() *Config {
	_ = godotenv.Load()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if cp := os.Getenv("CONFIG_PATH"); cp != "" {
		v.AddConfigPath(cp)
	}
	v.AddConfigPath("conf")
	v.AddConfigPath("./conf")
	v.AddConfigPath("../conf")
	v.AddConfigPath("../../conf")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config: read file error: %v", err)
	}

	cfg := &Config{
		PostgresHost: v.GetString("postgres.host"),
		PostgresPort: v.GetInt("postgres.port"),
		PostgresDB:   v.GetString("postgres.db"),
	}

	cfg.PostgresUser = os.Getenv("POSTGRES_USER")
	cfg.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")

	if cfg.PostgresHost == "" || cfg.PostgresPort == 0 || cfg.PostgresDB == "" {
		log.Fatal("config: postgres.host, postgres.port, postgres.db are required")
	}
	if cfg.PostgresUser == "" || cfg.PostgresPassword == "" {
		log.Fatal("config: POSTGRES_USER and POSTGRES_PASSWORD env vars are required")
	}

	return cfg
}
