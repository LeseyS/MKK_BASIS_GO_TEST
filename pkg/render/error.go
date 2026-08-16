package render

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/logger"
)

type Err struct {
	Error string `json:"error"`
}

func Error(ctx context.Context, w http.ResponseWriter, err error, status int, message string) {
	ctxErr, ok := ctx.Value(logger.ContextErrKey{}).(*error)
	if ok {
		*ctxErr = fmt.Errorf("%s: %w", message, err)
	}

	err = unpack(err)
	err = fmt.Errorf("%s: %w", message, err)

	JSON(w, Err{Error: err.Error()}, status)
}

func unpack(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			break
		}

		err = e
	}

	return err
}
