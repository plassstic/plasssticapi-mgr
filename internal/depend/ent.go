package depend

import (
	"context"
	"fmt"

	"plassstic.tech/gopkg/golang-manager/lib/ent"

	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func NewEntClient(lc fx.Lifecycle, config *Config, log *zap.SugaredLogger) *ent.Client {
	client, err := ent.Open("postgres", config.PostgresData)

	if err != nil {
		log.Named("database.NewEntClient").Panic(fmt.Sprintf("panic! <%T> %v", err, err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
