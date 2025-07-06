package storage

import (
	"auth_service/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
)

func Migrate(cfg *config.Config) error {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Dbname,
	)

	const delay = 5 * time.Second

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// временное подключение к бд для миграций
	for i := 0; i < cfg.Database.MaxAttempts; i++ {

		// проверка соединенияы
		if err := db.PingContext(ctxTimeOut); err == nil {
			break
		}

		// миграции
		if i == cfg.Database.MaxAttempts-1 {
			return errors.New("не удалось подключиться к базе данных для миграций")
		}
		time.Sleep(delay)

	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	// миграции
	if err := goose.Up(db, os.Getenv("MIGRATIONS_DIR")); err != nil {
		return errors.New("ошибка миграций" + err.Error())
	}

	return nil
}
