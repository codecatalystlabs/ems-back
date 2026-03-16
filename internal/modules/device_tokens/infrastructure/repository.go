package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	deviceapp "dispatch/internal/modules/device_tokens/application"
	devicedomain "dispatch/internal/modules/device_tokens/domain"
	platformdb "dispatch/internal/platform/db"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var _ deviceapp.Repository = (*Repository)(nil)

func (r *Repository) Register(ctx context.Context, in devicedomain.DeviceToken) (devicedomain.DeviceToken, error) {
	now := time.Now().UTC()

	// if push token exists already, reactivate/update instead of failing
	err := r.db.QueryRow(ctx, `
		INSERT INTO user_device_tokens (
			id, user_id, device_id, platform, push_token, is_active, last_seen_at, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,TRUE,$6,$7,$7)
		ON CONFLICT (push_token) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			device_id = EXCLUDED.device_id,
			platform = EXCLUDED.platform,
			is_active = TRUE,
			last_seen_at = EXCLUDED.last_seen_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, user_id, COALESCE(device_id,''), platform, push_token, is_active, last_seen_at, created_at, updated_at
	`,
		in.ID, in.UserID, nullIfEmpty(in.DeviceID), in.Platform, in.PushToken, now, now,
	).Scan(
		&in.ID, &in.UserID, &in.DeviceID, &in.Platform, &in.PushToken,
		&in.IsActive, &in.LastSeenAt, &in.CreatedAt, &in.UpdatedAt,
	)

	return in, err
}

func (r *Repository) Update(ctx context.Context, id string, req deviceapp.UpdateDeviceTokenRequest) (devicedomain.DeviceToken, error) {
	sets := make([]string, 0)
	args := make([]any, 0)
	pos := 1

	if req.DeviceID != nil {
		sets = append(sets, fmt.Sprintf("device_id = $%d", pos))
		args = append(args, nullIfEmpty(*req.DeviceID))
		pos++
	}
	if req.Platform != nil {
		sets = append(sets, fmt.Sprintf("platform = $%d", pos))
		args = append(args, *req.Platform)
		pos++
	}
	if req.PushToken != nil {
		sets = append(sets, fmt.Sprintf("push_token = $%d", pos))
		args = append(args, *req.PushToken)
		pos++
	}
	if req.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", pos))
		args = append(args, *req.IsActive)
		pos++
	}

	if len(sets) == 0 {
		return r.GetByID(ctx, id)
	}

	sets = append(sets, "updated_at = now()", "last_seen_at = now()")
	args = append(args, id)

	q := fmt.Sprintf(`
		UPDATE user_device_tokens
		SET %s
		WHERE id = $%d
	`, strings.Join(sets, ", "), pos)

	ct, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return devicedomain.DeviceToken{}, err
	}
	if ct.RowsAffected() == 0 {
		return devicedomain.DeviceToken{}, deviceapp.ErrDeviceTokenNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) Deactivate(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE user_device_tokens
		SET is_active = FALSE,
		    updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return deviceapp.ErrDeviceTokenNotFound
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (devicedomain.DeviceToken, error) {
	var out devicedomain.DeviceToken
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, COALESCE(device_id,''), platform, push_token, is_active, last_seen_at, created_at, updated_at
		FROM user_device_tokens
		WHERE id = $1
	`, id).Scan(
		&out.ID, &out.UserID, &out.DeviceID, &out.Platform, &out.PushToken,
		&out.IsActive, &out.LastSeenAt, &out.CreatedAt, &out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return devicedomain.DeviceToken{}, deviceapp.ErrDeviceTokenNotFound
	}
	return out, err
}

func (r *Repository) List(ctx context.Context, params deviceapp.ListDeviceTokensParams) ([]devicedomain.DeviceToken, int64, error) {
	p := params.Pagination

	where := []string{"1=1"}
	args := make([]any, 0)
	pos := 1

	if params.UserID != nil && *params.UserID != "" {
		where = append(where, fmt.Sprintf("udt.user_id = $%d", pos))
		args = append(args, *params.UserID)
		pos++
	}

	if params.Platform != nil && *params.Platform != "" {
		where = append(where, fmt.Sprintf("udt.platform = $%d", pos))
		args = append(args, strings.ToUpper(*params.Platform))
		pos++
	}

	if params.IsActive != nil {
		where = append(where, fmt.Sprintf("udt.is_active = $%d", pos))
		args = append(args, *params.IsActive)
		pos++
	}

	if p.Search != "" {
		where = append(where, fmt.Sprintf("(COALESCE(udt.device_id,'') ILIKE $%d OR udt.push_token ILIKE $%d)", pos, pos))
		args = append(args, "%"+p.Search+"%")
		pos++
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(1) FROM user_device_tokens udt `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := fmt.Sprintf(`
		SELECT id, user_id, COALESCE(device_id,''), platform, push_token, is_active, last_seen_at, created_at, updated_at
		FROM user_device_tokens udt
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereSQL, platformdb.BuildOrderBy(p, map[string]string{
		"created_at":   "udt.created_at",
		"updated_at":   "udt.updated_at",
		"platform":     "udt.platform",
		"last_seen_at": "udt.last_seen_at",
	}), pos, pos+1)

	rows, err := r.db.Query(ctx, q, append(args, p.PageSize, p.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []devicedomain.DeviceToken{}
	for rows.Next() {
		var out devicedomain.DeviceToken
		if err := rows.Scan(
			&out.ID, &out.UserID, &out.DeviceID, &out.Platform, &out.PushToken,
			&out.IsActive, &out.LastSeenAt, &out.CreatedAt, &out.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, out)
	}
	return items, total, rows.Err()
}

func (r *Repository) GetPushTokensByUserID(ctx context.Context, userID string) ([]string, error) {
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

func nullIfEmpty(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}
