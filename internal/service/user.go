package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/repository"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent/schema"
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
	srv.result = srv.
		repo.
		GetByID(ctx, userId).
		Result()

	return srv
}

func (srv *UserService) GetAll() *UserService {
	srv.result = srv.
		repo.
		GetAll(ctx).
		Result()

	return srv
}

func (srv *UserService) FilterMany(filter func(item *ent.User, index int) bool) *UserService {
	var slice []*ent.User
	var ok bool

	if slice, ok = srv.result.([]*ent.User); !ok {
		return srv
	}

	srv.result = lo.Filter(slice, filter)

	return srv
}

func (srv *UserService) Ensure(userId int) *UserService {
	srv.result = srv.
		repo.
		GetByID(ctx, userId).
		ClearIfErrAs(&ent.NotFoundError{}).
		Create(ctx, userId).
		Result()

	return srv
}

func (srv *UserService) SetBoth(userId int64, bot schema.Bot, editable schema.Editable) *UserService {
	srv.result = srv.
		repo.
		SetBoth(ctx, userId, bot, editable).
		Result()

	return srv
}

func (srv *UserService) One() (*ent.User, error) {
	switch srv.result.(type) {
	case *ent.User:
		_ = srv.tx.Commit()
		return srv.result.(*ent.User), nil
	case error:
		_ = srv.tx.Rollback()
		return nil, srv.result.(error)
	default:
		return nil, fmt.Errorf("result is empty")
	}
}

func (srv *UserService) Many() ([]*ent.User, error) {
	switch srv.result.(type) {
	case []*ent.User:
		_ = srv.tx.Commit()
		return srv.result.([]*ent.User), nil
	case error:
		_ = srv.tx.Rollback()
		return nil, srv.result.(error)
	default:
		return nil, fmt.Errorf("result is empty")
	}
}
