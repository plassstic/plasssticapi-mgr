package depend

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/api"
)

type Config struct {
	PostgresData string
	BotToken     string
}

func NewConfig() *Config {
	log := GetLogger("depend.spawnConfig")
	data := getPostgresData()
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Panic("panic! environment variable BOT_TOKEN is unspecified")
	}

	apiDomain := os.Getenv("API_DOMAIN")

	if apiDomain == "" {
		log.Panic("panic! environment variable API_DOMAIN is unspecified")
	}

	Domain = apiDomain

	return &Config{
		PostgresData: data,
		BotToken:     botToken,
	}
}

func getPostgresData() string {
	log := GetLogger("config.getPostgresData")

	if err := godotenv.Load(); err != nil {
		log.Panicf("panic! <%T> %e", err, err)
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
