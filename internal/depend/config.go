package depend

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	PostgresData string
	BotToken     string
}

func NewConfig(log *zap.SugaredLogger) *Config {
	data := getPostgresData(log)
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Named("config.NewConfig").Panic("panic! environment variable BOT_TOKEN is unspecified")
	}
	return &Config{
		PostgresData: data,
		BotToken:     botToken,
	}
}

func getPostgresData(log *zap.SugaredLogger) string {
	err := godotenv.Load()

	if err != nil {
		log.Named("config.getPostgresData").Panic(fmt.Sprintf("panic! <%T> %v", err, err))
	}

	url := fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	return url
}
