package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	Db         DbConfig
	Mail       MailConfig
	GoogleJson string
}

type DbConfig struct {
	User   string
	Pass   string
	Port   string
	Server string
	DbName string
}

type MailConfig struct {
	User        string
	Pass        string
	Host        string
	Port        string
	From        string
	Admin       string
	SendGridKey string
}

var config = Config{}

func InitConfig(skipFile bool) error {
	if !skipFile {
		env := os.Getenv("GO_ENV")
		if env == "" {
			env = "dev"
		}
		fmt.Printf("env:%s\n", env)
		err := godotenv.Load(fmt.Sprintf("./.env.%s", env))
		if err != nil {
			return fmt.Errorf("failed to load env. %v", err)
		}
	}

	config = Config{}
	config.AppPort = os.Getenv("APP_PORT")
	if config.AppPort == "" {
		return errors.New("no AppPort")
	}
	config.GoogleJson = os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if config.GoogleJson == "" {
		return errors.New("no Google Credendials")
	}

	config.Db = newDbConfig()
	config.Mail = newMailConfig()

	return nil
}

func GetConfig() Config {
	return config
}

func newDbConfig() DbConfig {
	config := DbConfig{
		User:   os.Getenv("DB_USER"),
		Pass:   os.Getenv("DB_PASS"),
		Port:   os.Getenv("DB_PORT"),
		Server: os.Getenv("DB_SERVER"),
		DbName: os.Getenv("DB_NAME"),
	}
	return config
}

func newMailConfig() MailConfig {
	config := MailConfig{
		User:        os.Getenv("MAIL_USER"),
		Pass:        os.Getenv("MAIL_PASS"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        os.Getenv("MAIL_PORT"),
		From:        os.Getenv("MAIL_FROM"),
		Admin:       os.Getenv("MAIL_ADMIN"),
		SendGridKey: os.Getenv("SEND_GRID_API_KEY"),
	}
	return config
}
