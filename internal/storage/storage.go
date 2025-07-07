package storage

import (
	"auth_service/internal/config"
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type RefreshSession struct {
	ID         int
	UserID     uuid.UUID
	TokenHash  string
	UserAgent  string
	IPAddress  string
	ExpiryDate time.Time
	IsRevoked  bool
}

type Storage interface {
	CreateRefreshSession(ctx context.Context, session *RefreshSession) error
	FindRefreshSessions(ctx context.Context, userID uuid.UUID) ([]RefreshSession, error)
	DeleteRefreshSession(ctx context.Context, sessionID int) error
	RevokeUserSessions(ctx context.Context, userID uuid.UUID) error
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) error
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

func (s *Store) CreateRefreshSession(ctx context.Context, session *RefreshSession) error {

	query := `
	INSERT INTO refresh_sessions (user_id, token_hash, user_agent, ip_address, expiry_date, is_revoked)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	var id int
	err := s.Pool.QueryRow(ctx, query, session.UserID, session.TokenHash, session.UserAgent, session.IPAddress, session.ExpiryDate, session.IsRevoked).Scan(&id)

	return err
}

func (s *Store) FindRefreshSessions(ctx context.Context, userID uuid.UUID) ([]*RefreshSession, error) {
	query := `
		SELECT * FROM refresh_sessions
		WHERE user_id = $1 AND expiry_date > NOW() AND is_revoked = FALSE
	`
	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*RefreshSession
	for rows.Next() {
		var s RefreshSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.TokenHash, &s.UserAgent, &s.IPAddress, &s.ExpiryDate, &s.IsRevoked); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (s *Store) DeleteRefreshSession(ctx context.Context, sessionID int) error {
	query := "UPDATE refresh_sessions SET is_revoked = TRUE WHERE id = $1"

	_, err := s.Pool.Exec(ctx, query, sessionID)
	return err
}

func (s *Store) RevokeUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := "UPDATE refresh_sessions SET is_revoked = TRUE WHERE user_id = $1 AND is_revoked = FALSE"

	_, err := s.Pool.Exec(ctx, query, userID)
	return err
}

func (s *Store) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM refresh_sessions WHERE user_id = $1`

	var count int
	err := s.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
