package unitofwork

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

//registra todos repositorios
type UnitOfWorkInterface interface {
	Register(name string, fc RepositoryFactory) // registrar o repositorio
	GetRepository(ctx context.Context, name string) (interface{}, error) // pegar o repositorio
	Do(ctx context.Context, fn func(UnitOfWork *UnitOfWork) error) error //realiza a acao
	CommitOrRollback() error //fazer o commit da acao ou rollback caso der erro em alguma operacao
	Rollback() error
	UnRegister(name string)
}

type UnitOfWork struct {
	Db           *sql.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

func NewUnitOfWork(ctx context.Context, db *sql.DB) *UnitOfWork {
	return &UnitOfWork{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *UnitOfWork) Register(name string, fc RepositoryFactory) {
	u.Repositories[name] = fc
}

func (u *UnitOfWork) UnRegister(name string) {
	delete(u.Repositories, name)
}

func (u *UnitOfWork) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		tx, err := u.Db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.Tx = tx
	}
	repo := u.Repositories[name](u.Tx)
	return repo, nil
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(UnitOfWork *UnitOfWork) error) error {
	if u.Tx != nil {
		return fmt.Errorf("transaction already started")
	}
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx
	err = fn(u)
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return errors.New(fmt.Sprintf("original error: %s, rollback error: %s", err.Error(), errRb.Error()))
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UnitOfWork) Rollback() error {
	if u.Tx == nil {
		return errors.New("no transaction to rollback")
	}
	err := u.Tx.Rollback()
	if err != nil {
		return err
	}
	u.Tx = nil
	return nil
}

func (u *UnitOfWork) CommitOrRollback() error {
	err := u.Tx.Commit()
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return errors.New(fmt.Sprintf("original error: %s, rollback error: %s", err.Error(), errRb.Error()))
		}
		return err
	}
	u.Tx = nil
	return nil
}