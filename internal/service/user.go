package service

import (
	"context"

	. "plassstic.tech/gopkg/golang-manager/internal/repository"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

var ctx = context.Background()

type UserService struct {
	repo   *UserRepo
	tx     *ent.Tx
	result interface{}
	err    error
}

func (srv *UserService) With(client *ent.Client) *UserService {
	srv.tx, _ = client.Tx(context.Background())
	srv.repo = (&UserRepo{}).With(srv.tx)
	return srv
}

func (srv *UserService) Get(userId int) *UserService {
	srv.result = srv.repo.
		GetByID(ctx, userId).
		Result()

	return srv
}

func (srv *UserService) Ensure(userId int) *UserService {
	srv.result = srv.repo.
		GetByID(ctx, userId).
		ClearIfErrAs(&ent.NotFoundError{}).
		Create(ctx, userId).
		Result()

	return srv
}

func (srv *UserService) Fin() interface{} {
	switch srv.result.(type) {
	case error:
		_ = srv.tx.Rollback()
	default:
		_ = srv.tx.Commit()
	}

	//logger.GetLogger("srv").Debugf("Return %#v", srv.result)

	return srv.result
}
