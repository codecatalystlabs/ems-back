package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dispatch/internal/modules/users/application/dto"
	"dispatch/internal/modules/users/domain"
	"dispatch/internal/platform/db"

	userapp "dispatch/internal/modules/users/application"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u domain.User, passwordHash string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, username, first_name, last_name, phone, email, password_hash, status, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,true,$9,$9)
	`, u.ID, u.Username, u.FirstName, u.LastName, u.Phone, u.Email, passwordHash, u.Status, u.CreatedAt)
	return err
}

func (r *Repository) List(ctx context.Context, params dto.ListUsersParams) ([]domain.User, int64, error) {
	p := params.Pagination
	allowedSorts := map[string]string{
		"created_at": "u.created_at",
		"username":   "u.username",
		"first_name": "u.first_name",
		"last_name":  "u.last_name",
		"status":     "u.status",
	}

	baseWhere := []string{"u.deleted_at IS NULL"}
	args := make([]any, 0)
	argPos := 1

	if p.Search != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("(u.username ILIKE $%d OR u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR COALESCE(u.email, '') ILIKE $%d OR COALESCE(u.phone, '') ILIKE $%d)", argPos, argPos, argPos, argPos, argPos))
		args = append(args, "%"+p.Search+"%")
		argPos++
	}

	for key, value := range p.Filters {
		switch key {
		case "status":
			baseWhere = append(baseWhere, fmt.Sprintf("u.status = $%d", argPos))
			args = append(args, strings.ToUpper(value))
			argPos++
		case "is_active":
			baseWhere = append(baseWhere, fmt.Sprintf("u.is_active::text = $%d", argPos))
			args = append(args, strings.ToLower(value))
			argPos++
		}
	}

	whereSQL := "WHERE " + strings.Join(baseWhere, " AND ")
	countSQL := fmt.Sprintf(`SELECT COUNT(1) FROM users u %s`, whereSQL)

	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := db.BuildOrderBy(p, allowedSorts)
	listSQL := fmt.Sprintf(`
		SELECT u.id, u.username, u.first_name, u.last_name, COALESCE(u.phone,''), COALESCE(u.email,''), u.status, u.created_at
		FROM users u
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argPos, argPos+1)

	listArgs := append(args, p.PageSize, p.Offset)
	rows, err := r.db.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Phone, &u.Email, &u.Status, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, username, first_name, last_name, COALESCE(phone,''), COALESCE(email,''), status, created_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Phone, &u.Email, &u.Status, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, userapp.ErrUserNotFound
	}
	return u, err
}

func (r *Repository) Update(ctx context.Context, id string, req dto.UpdateUserRequest) (domain.User, error) {
	sets := make([]string, 0)
	args := make([]any, 0)
	pos := 1

	if req.FirstName != nil {
		sets = append(sets, fmt.Sprintf("first_name = $%d", pos))
		args = append(args, *req.FirstName)
		pos++
	}
	if req.LastName != nil {
		sets = append(sets, fmt.Sprintf("last_name = $%d", pos))
		args = append(args, *req.LastName)
		pos++
	}
	if req.Phone != nil {
		sets = append(sets, fmt.Sprintf("phone = $%d", pos))
		args = append(args, *req.Phone)
		pos++
	}
	if req.Email != nil {
		sets = append(sets, fmt.Sprintf("email = $%d", pos))
		args = append(args, *req.Email)
		pos++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", pos))
		args = append(args, strings.ToUpper(*req.Status))
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

	sets = append(sets, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
	`, strings.Join(sets, ", "), pos)

	ct, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return domain.User{}, err
		}
		return domain.User{}, err
	}
	if ct.RowsAffected() == 0 {
		return domain.User{}, userapp.ErrUserNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE users
		SET deleted_at = now(),
		    is_active = FALSE,
		    updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return userapp.ErrUserNotFound
	}
	return nil
}
