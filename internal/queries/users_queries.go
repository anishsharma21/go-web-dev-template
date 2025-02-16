package queries

import (
	"context"
	"fmt"

	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserByEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"

	row := dbPool.QueryRow(ctx, query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return models.User{}, fmt.Errorf("error retrieving user with email %q: %w", email, err)
	}
	return user, nil
}
