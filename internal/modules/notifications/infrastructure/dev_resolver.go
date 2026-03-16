package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceTokenRepository struct {
	db *pgxpool.Pool
}

func NewDeviceTokenRepository(db *pgxpool.Pool) *DeviceTokenRepository {
	return &DeviceTokenRepository{db: db}
}

func (r *DeviceTokenRepository) GetPushTokensByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT push_token
		FROM user_device_tokens
		WHERE user_id = $1
		  AND is_active = TRUE
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}
