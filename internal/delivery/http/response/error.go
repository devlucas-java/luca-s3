package response

import (
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/domain/errors"
)

func ResponseError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError
	switch e := err.(type) {
	case *errors.AppError:
		appErr = e
	default:
		appErr = errors.New("INTERNAL_ERROR", "an unexpected error occurred", http.StatusInternalServerError, err)
	}

	ResponseEntity(w, appErr.Status, appErr)
}
