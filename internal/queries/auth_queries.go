package queries

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SignUpNewUser(ctx context.Context, dbPool *pgxpool.Pool, user models.User) error {
	args := pgx.NamedArgs{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"password":   user.Password,
	}

	query := "INSERT INTO users (email, first_name, last_name, password) VALUES (@email, @first_name, @last_name, @password)"

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if ct.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to insert user: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("User signed up successfully: %s", user.Email), "command_tag", ct.String())

	return nil
}
