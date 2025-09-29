package repository

import (
	"context"
	"errors"

	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
	"plassstic.tech/gopkg/golang-manager/lib/ent/schema"
)

type UserRepo struct {
	tx     *ent.Tx
	result interface{}
}

func (repo *UserRepo) With(tx *ent.Tx) *UserRepo {
	repo.tx = tx
	return repo
}
func (repo *UserRepo) SetBoth(ctx context.Context, userId int64, bot schema.Bot, editable schema.Editable) *UserRepo {
	logger.GetLogger("123").Infof("%#v", repo)

	if repo.resultNotNil() {
		return repo
	}

	var e error

	repo.result, e = repo.
		tx.
		Client().
		User.
		UpdateOneID(int(userId)).
		SetBot(bot).
		SetEditable(editable).
		Save(ctx)

	if e != nil {
		repo.result = e
	}

	return repo
}
func (repo *UserRepo) GetByID(ctx context.Context, userId int) *UserRepo {
	if repo.resultNotNil() {
		return repo
	}

	var e error

	repo.result, e = repo.
		tx.
		Client().
		User.
		Get(ctx, userId)

	//logger.GetLogger("repo").Debugf("Got %v (err %v)", repo.result, e)

	if e != nil {
		repo.result = e
	}

	return repo
}

func (repo *UserRepo) Create(ctx context.Context, userId int) *UserRepo {
	if repo.resultNotNil() {
		return repo
	}

	var e error

	repo.result, e = repo.
		tx.
		Client().
		User.
		Create().
		SetID(userId).
		Save(ctx)

	//logger.GetLogger("repo").Debugf("Created %#v (err %#v)", repo.result, e)

	if e != nil {
		repo.result = e
	}

	return repo
}

func (repo *UserRepo) TransformResult(callable func(interface{}) interface{}) *UserRepo {
	repo.result = callable(repo.result)
	return repo
}

func (repo *UserRepo) ClearIfErrIs(err error) *UserRepo {
	if repo.compareResultErrIs(err) {
		//logger.GetLogger("repo").Debugf("Result is %#v", err)
		return repo.Clear()
	}

	return repo
}

func (repo *UserRepo) ClearIfErrAs(err interface{}) *UserRepo {
	if repo.compareResultErrAs(&err) {
		//logger.GetLogger("repo").Debugf("Result is %#v", err)
		return repo.Clear()
	}

	return repo
}

func (repo *UserRepo) Clear() *UserRepo {
	repo.result = nil
	return repo
}

func (repo *UserRepo) resultIsErr() bool {
	_, ok := repo.result.(error)
	return ok
}

func (repo *UserRepo) resultNotNil() bool {
	return repo.result != nil
}

func (repo *UserRepo) TResultErrFunc(f func(error) error) *UserRepo {
	if repo.resultIsErr() {
		repo.result = f(repo.result.(error))
	}

	return repo
}
func (repo *UserRepo) compareResultErrAs(err *interface{}) bool {
	if repo.resultIsErr() {
		return errors.As(repo.result.(error), err)
	}
	return false
}
func (repo *UserRepo) compareResultErrIs(err error) bool {
	if repo.resultIsErr() {
		return errors.Is(repo.result.(error), err)
	}
	return false
}

func (repo *UserRepo) Result() interface{} {
	//logger.GetLogger("repo").Debugf("ret %#v", repo.result)
	return repo.result
}
