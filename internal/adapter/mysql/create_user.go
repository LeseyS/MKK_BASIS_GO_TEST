package mysql

import (
	"context"
	"fmt"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/apperr"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/internal/domain"
	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/transaction"
)

func (my *MySQL) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	const query = `INSERT INTO users (email, name, password_hash) VALUES (?, ?, ?)`

	args := []any{
		user.Email,
		user.Name,
		user.PasswordHash,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	res, err := txOrPool.ExecContext(ctx, query, args...)
	if err != nil {
		if isDuplicateErr(err) {
			return 0, apperr.ErrDuplicate
		}
		return 0, fmt.Errorf("inserting user: %w", err)
	}

	return res.LastInsertId()
}
