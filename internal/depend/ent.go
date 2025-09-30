package depend

import (
	"context"

	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func NewEntClient(lc fx.Lifecycle, config *Config) *ent.Client {
	var client *ent.Client
	var err error

	log := GetLogger("depend.spawnEnt")

	if client, err = ent.Open(
		"postgres",
		config.PostgresData,
		ent.Log(GetLogger("ent").Info),
	); err != nil {
		log.Panicf("panic! <%T> %e", err, err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
