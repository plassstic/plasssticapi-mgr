package depend

import (
	"context"
	"fmt"

	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	"plassstic.tech/gopkg/golang-manager/lib/ent"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func NewEntClient(lc fx.Lifecycle, config *Config) *ent.Client {
	log := logger.GetLogger("database.NewEntClient")
	client, err := ent.Open("postgres", config.PostgresData)

	if err != nil {
		log.Panic(fmt.Sprintf("panic! <%T> %v", err, err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
