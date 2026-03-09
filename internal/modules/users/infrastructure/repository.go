package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"dispatch/internal/modules/users/application/dto"
	"dispatch/internal/modules/users/domain"
	"dispatch/internal/platform/db"

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
