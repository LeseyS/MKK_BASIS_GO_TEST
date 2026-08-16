package mysql

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNoRows = sql.ErrNoRows
)

func isDuplicateErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
