package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `SELECT id, email, name, password_hash, created_at FROM users WHERE email = ?`

	txOrPool := transaction.TryExtractTX(ctx)

	var user domain.User

	err := txOrPool.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return user, apperr.ErrNotFound
	}

	if err != nil {
		return user, fmt.Errorf("txOrPool.QueryRowContext: %w", err)
	}

	return user, nil
}
