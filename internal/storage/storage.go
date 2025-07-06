package storage

import (
	"auth_service/internal/config"
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Storage interface {
}

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Open(ctx context.Context, cfg *config.Config, logger *logrus.Logger) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Dbname,
	)

	if err := Migrate(cfg); err != nil {
		return err
	}

	delay := time.Second

	ctxTimeOut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for i := 0; i < cfg.Database.MaxAttempts; i++ {

		dbpool, err := pgxpool.New(ctxTimeOut, dsn)
		if err == nil {
			if err = dbpool.Ping(ctxTimeOut); err == nil {
				s.Pool = dbpool
				return nil
			} else {
				logger.Debug("ошибка пинга базы данных: ", err)
			}

			dbpool.Close()
		} else {
			logger.Debug("ошибка создания пула соединений: ", err)
		}

		select {
		case <-time.After(delay):
			logger.Debugf("попытка %d из %d подключения к базе данных", i+1, cfg.Database.MaxAttempts)
		case <-ctxTimeOut.Done():
			return errors.New("истекло время для подключения к базе данных")
		}
	}

	return errors.New("превышено количество попыток подключения к базе данных")
}

func (s *Store) Close() error {
	if s.Pool == nil {
		return errors.New("нет пула соединений")
	}
	s.Pool.Close()
	return nil
}
